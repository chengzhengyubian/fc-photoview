package scanner

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/scanner_cache"
	"github.com/photoview/photoview/api/scanner/scanner_tasks/cleanup_tasks"
	"github.com/photoview/photoview/api/scanner/scanner_utils"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
	ignore "github.com/sabhiram/go-gitignore"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func getPhotoviewIgnore(ignorePath string) ([]string, error) {
	var photoviewIgnore []string

	// Open .photoviewignore file, if exists
	photoviewIgnoreFile, err := os.Open(path.Join(ignorePath, ".photoviewignore"))
	if err != nil {
		if os.IsNotExist(err) {
			return photoviewIgnore, nil
		}
		return photoviewIgnore, err
	}

	// Close file on exit
	defer photoviewIgnoreFile.Close()

	// Read and save .photoviewignore data
	scanner := bufio.NewScanner(photoviewIgnoreFile)
	for scanner.Scan() {
		photoviewIgnore = append(photoviewIgnore, scanner.Text())
		log.Printf("Ignore found: %s", scanner.Text())
	}

	return photoviewIgnore, scanner.Err()
}

func FindAlbumsForUser(user *models.User, album_cache *scanner_cache.AlbumScannerCache) ([]*models.Album, []error) {

	if err := user.FillAlbums(); err != nil {
		return nil, []error{err}
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}

	var userRootAlbums []*models.Album
	userAlbumid, err := json.Marshal(userAlbumIDs)
	if err != nil {
		fmt.Print(err)
	}
	userAlbumids := strings.Trim(string(userAlbumid), "[]")
	sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE id IN (%v) AND (parent_album_id IS NULL OR parent_album_id NOT IN (%v))", userAlbumids, userAlbumids)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	num := len(res)
	if num == 0 {
		return nil, []error{err}
	}
	for i := 0; i < num; i++ {
		var userRootAlbum models.Album
		userRootAlbum.ID = int(*res[i][0].LongValue)
		userRootAlbum.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		userRootAlbum.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		userRootAlbum.Title = *res[i][3].StringValue
		if res[i][4].IsNull != nil {
			userRootAlbum.ParentAlbumID = nil
		} else {
			id := int(*res[i][4].LongValue)
			userRootAlbum.ParentAlbumID = &id
		}
		userRootAlbum.Path = *res[i][5].StringValue
		userRootAlbum.PathHash = *res[i][6].StringValue
		if res[i][7].IsNull != nil {
			userRootAlbum.CoverID = nil
		} else {
			id := int(*res[i][7].LongValue)
			userRootAlbum.CoverID = &id
		}
		userRootAlbums = append(userRootAlbums, &userRootAlbum)
	}
	scanErrors := make([]error, 0)

	type scanInfo struct {
		path   string
		parent *models.Album
		ignore []string
	}

	scanQueue := list.New()

	for _, album := range userRootAlbums {
		// Check if user album directory exists on the file system
		if _, err := os.Stat(album.Path); err != nil {
			if os.IsNotExist(err) {
				scanErrors = append(scanErrors, errors.Errorf("Album directory for user '%s' does not exist '%s'\n", user.Username, album.Path))
			} else {
				scanErrors = append(scanErrors, errors.Errorf("Could not read album directory for user '%s': %s\n", user.Username, album.Path))
			}
		} else {
			scanQueue.PushBack(scanInfo{
				path:   album.Path,
				parent: nil,
				ignore: nil,
			})
		}
	}

	userAlbums := make([]*models.Album, 0)

	for scanQueue.Front() != nil {
		albumInfo := scanQueue.Front().Value.(scanInfo)
		scanQueue.Remove(scanQueue.Front())

		albumPath := albumInfo.path
		albumParent := albumInfo.parent
		albumIgnore := albumInfo.ignore

		// Read path
		dirContent, err := ioutil.ReadDir(albumPath)
		if err != nil {
			scanErrors = append(scanErrors, errors.Wrapf(err, "read directory (%s)", albumPath))
			continue
		}

		// Skip this dir if in ignore list
		ignorePaths := ignore.CompileIgnoreLines(albumIgnore...)
		if ignorePaths.MatchesPath(albumPath + "/") {
			log.Printf("Skip, directroy %s is in ignore file", albumPath)
			continue
		}

		// Update ignore dir list
		photoviewIgnore, err := getPhotoviewIgnore(albumPath)
		if err != nil {
			log.Printf("Failed to get ignore file, err = %s", err)
		} else {
			albumIgnore = append(albumIgnore, photoviewIgnore...)
		}

		// Will become new album or album from db
		var album *models.Album

		//	transErr := db.Transaction(func(tx *gorm.DB) error {
		log.Printf("Scanning directory: %s", albumPath)

		// check if album already exists
		var albumResult []models.Album

		sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE path_hash = '%v'", models.MD5Hash(albumPath))
		res, err := dataApi.Query(sql_albums_se)
		if err != nil {
			fmt.Println(err)
		}
		num := len(res)
		for i := 0; i < num; i++ {
			var Album models.Album
			Album.ID = DataApi.GetInt(res, i, 0)
			Album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
			Album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
			Album.Title = DataApi.GetString(res, i, 3)
			Album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
			Album.Path = DataApi.GetString(res, i, 5)
			Album.PathHash = DataApi.GetString(res, i, 6)
			Album.CoverID = DataApi.GetIntP(res, i, 7)
			albumResult = append(albumResult, Album)
		}
		// album does not exist, create new
		if len(albumResult) == 0 {
			albumTitle := path.Base(albumPath)

			var albumParentID *int
			parentOwners := make([]models.User, 0)
			if albumParent != nil {
				albumParentID = &albumParent.ID

				sql_users_se := fmt.Sprintf("SELECT `users`.`id`,`users`.`created_at`,`users`.`updated_at`,`users`.`username`,`users`.`password`,`users`.`admin` FROM `users` JOIN `user_albums` ON `user_albums`.`user_id` = `users`.`id` AND `user_albums`.`album_id` = %v", albumParentID)
				res, err = dataApi.Query(sql_users_se)
				num = len(res)
				for i := 0; i < num; i++ {
					var User models.User
					User.ID = DataApi.GetInt(res, i, 0)
					//User.CreatedAt = time.Unix(DataApi.GetLong(res, i, 1)/1000, 0)
					//User.UpdatedAt = time.Unix(DataApi.GetLong(res, i, 2)/1000, 0)
					User.Username = DataApi.GetString(res, i, 3)
					User.Password = DataApi.GetStringP(res, i, 4)
					User.Admin = DataApi.GetBoolean(res, i, 5)
					parentOwners = append(parentOwners, User)
				}
			}

			album = &models.Album{
				Title:         albumTitle,
				ParentAlbumID: albumParentID,
				Path:          albumPath,
			}

			// Store album ignore
			album_cache.InsertAlbumIgnore(albumPath, albumIgnore)
			var ParentAlbumid int
			ParentAlbumid = *album.ParentAlbumID
			album.PathHash = models.MD5Hash(album.Path)
			sql_albums_in := fmt.Sprintf("INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`) VALUES (NOW(),NOW(),'%v',%v,'%v','%v',NULL)", album.Title, ParentAlbumid, album.Path, album.PathHash)
			dataApi.ExecuteSQl(sql_albums_in)

			n := len(parentOwners)
			for i := 0; i < n; i++ {
				var user models.User
				user = parentOwners[i]
				sql_users_in := fmt.Sprintf("INSERT INTO `users` (`created_at`,`updated_at`,`username`,`password`,`admin`,`id`) VALUES ('%v','%v','%v','%v',%v,%v) ON DUPLICATE KEY UPDATE `id`=id", user.CreatedAt, user.UpdatedAt, user.Username, user.Password, user.Admin, user.ID)
				sql_users_albums_in := fmt.Sprintf("INSERT INTO `user_albums` (`album_id`,`user_id`) VALUES (%v,%v) ON DUPLICATE KEY UPDATE `album_id`=`album_id`", album.ID, user.ID)
				dataApi.ExecuteSQl(sql_users_in)
				dataApi.ExecuteSQl(sql_users_albums_in)
			}
			sql_albums_up := fmt.Sprintf("UPDATE `albums` SET `updated_at`=NOW() WHERE `id` = %v", album.ID)
			dataApi.ExecuteSQl(sql_albums_up)
		} else {
			album = &albumResult[0]

			// Add user as an owner of the album if not already
			var userAlbumOwner []models.User

			sql_users_se := fmt.Sprintf("SELECT `users`.`id`,`users`.`created_at`,`users`.`updated_at`,`users`.`username`,`users`.`password`,`users`.`admin` FROM `users` JOIN `user_albums` ON `user_albums`.`user_id` = `users`.`id` AND `user_albums`.`album_id` = %v WHERE user_albums.user_id = %v", album.ID, user.ID)
			res, err = dataApi.Query(sql_users_se)
			num = len(res)
			for i := 0; i < num; i++ {
				var User models.User
				User.ID = DataApi.GetInt(res, i, 0)
				//User.CreatedAt = time.Unix(DataApi.GetLong(res, i, 1)/1000, 0)
				//User.UpdatedAt = time.Unix(DataApi.GetLong(res, i, 2)/1000, 0)
				User.Username = DataApi.GetString(res, i, 3)
				User.Password = DataApi.GetStringP(res, i, 4)
				User.Admin = DataApi.GetBoolean(res, i, 5)
				userAlbumOwner = append(userAlbumOwner, User)
			}

			if len(userAlbumOwner) == 0 {
				newUser := models.User{}
				newUser.ID = user.ID
				sql_users_albums_in := fmt.Sprintf("INSERT INTO `user_albums` (`album_id`,`user_id`) VALUES (%v,%v) ON DUPLICATE KEY UPDATE `album_id`=`album_id`", album.ID, newUser.ID)
				dataApi.ExecuteSQl(sql_users_albums_in)
			}

			// Update album ignore
			album_cache.InsertAlbumIgnore(albumPath, albumIgnore)
		}

		userAlbums = append(userAlbums, album)
		// Scan for sub-albums
		for _, item := range dirContent {
			subalbumPath := path.Join(albumPath, item.Name())

			// Skip if directory is hidden
			if path.Base(subalbumPath)[0:1] == "." {
				continue
			}

			isDirSymlink, err := utils.IsDirSymlink(subalbumPath)
			if err != nil {
				scanErrors = append(scanErrors, errors.Wrapf(err, "could not check for symlink target of %s", subalbumPath))
				continue
			}

			if (item.IsDir() || isDirSymlink) && directoryContainsPhotos(subalbumPath, album_cache, albumIgnore) {
				scanQueue.PushBack(scanInfo{
					path:   subalbumPath,
					parent: album,
					ignore: albumIgnore,
				})
			}
		}
	}
	DataApi.Stresstest()
	deleteErrors := cleanup_tasks.DeleteOldUserAlbums( /*db, */ userAlbums, user)
	scanErrors = append(scanErrors, deleteErrors...)

	return userAlbums, scanErrors
}

func directoryContainsPhotos(rootPath string, cache *scanner_cache.AlbumScannerCache, albumIgnore []string) bool {

	if contains_image := cache.AlbumContainsPhotos(rootPath); contains_image != nil {
		return *contains_image
	}

	scanQueue := list.New()
	scanQueue.PushBack(rootPath)

	scanned_directories := make([]string, 0)

	for scanQueue.Front() != nil {

		dirPath := scanQueue.Front().Value.(string)
		scanQueue.Remove(scanQueue.Front())

		scanned_directories = append(scanned_directories, dirPath)

		// Update ignore dir list
		photoviewIgnore, err := getPhotoviewIgnore(dirPath)
		if err != nil {
			log.Printf("Failed to get ignore file, err = %s", err)
		} else {
			albumIgnore = append(albumIgnore, photoviewIgnore...)
		}
		ignoreEntries := ignore.CompileIgnoreLines(albumIgnore...)

		dirContent, err := ioutil.ReadDir(dirPath)
		if err != nil {
			scanner_utils.ScannerError("Could not read directory (%s): %s\n", dirPath, err.Error())
			return false
		}

		for _, fileInfo := range dirContent {
			filePath := path.Join(dirPath, fileInfo.Name())

			isDirSymlink, err := utils.IsDirSymlink(filePath)
			if err != nil {
				log.Printf("Cannot detect whether %s is symlink to a directory. Pretending it is not", filePath)
				isDirSymlink = false
			}

			if fileInfo.IsDir() || isDirSymlink {
				scanQueue.PushBack(filePath)
			} else {
				if cache.IsPathMedia(filePath) {
					if ignoreEntries.MatchesPath(fileInfo.Name()) {
						log.Printf("Match found %s, continue search for media", fileInfo.Name())
						continue
					}
					log.Printf("Insert Album %s %s, contains photo is true", dirPath, rootPath)
					cache.InsertAlbumPaths(dirPath, rootPath, true)
					return true
				}
			}
		}

	}

	for _, scanned_path := range scanned_directories {
		log.Printf("Insert Album %s, contains photo is false", scanned_path)
		cache.InsertAlbumPath(scanned_path, false)
	}
	return false
}
