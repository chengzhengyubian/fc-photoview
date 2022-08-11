package scanner

/*修改完*/
import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"time"

	//DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/scanner/scanner_tasks"
	"github.com/photoview/photoview/api/scanner/scanner_utils"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	//"time"
)

//修改中
func NewRootAlbum(rootPath string, owner *models.User) (*models.Album, error) {

	if !ValidRootPath(rootPath) {
		return nil, ErrorInvalidRootPath
	}

	if !path.IsAbs(rootPath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		rootPath = path.Join(wd, rootPath)
	}

	owners := []models.User{
		*owner,
	}

	var matchedAlbums []models.Album
	//if err := db.Where("path_hash = ?", models.MD5Hash(rootPath)).Find(&matchedAlbums).Error; err != nil { //SELECT * FROM `albums` WHERE path_hash = '1f7d46bead15a13c555e549e19624753'
	//	return nil, err
	//}
	//sql_albums_se := "SELECT * FROM `albums` WHERE path_hash =\"" + models.MD5Hash(rootPath) + "\""
	sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE path_hash ='%v'", models.MD5Hash(rootPath))
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
	}
	num := len(res)
	for i := 0; i < num; i++ {
		var album models.Album
		album.ID = DataApi.GetInt(res, i, 0)
		album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, i, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
		album.Path = DataApi.GetString(res, i, 5)
		album.PathHash = DataApi.GetString(res, i, 6)
		album.CoverID = DataApi.GetIntP(res, i, 7)
		matchedAlbums = append(matchedAlbums, album)
	}
	if len(matchedAlbums) > 0 {
		album := matchedAlbums[0]

		var matchedUserAlbumCount int64
		//if err := db.Table("user_albums").Where("user_id = ?", owner.ID).Where("album_id = ?", album.ID).Count(&matchedUserAlbumCount).Error; err != nil { //SELECT count(*) FROM `user_albums` WHERE user_id = 13 AND album_id = 108
		//	return nil, err
		//}
		//sql_user_albums_se := "SELECT count(*) FROM `user_albums` WHERE user_id =" + strconv.Itoa(owner.ID) + " AND album_id =" + strconv.Itoa(album.ID)
		sql_user_albums_se := fmt.Sprintf("SELECT count(*) FROM `user_albums` WHERE user_id =%v and album_id = %v", owner.ID, album.ID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err = dataApi.Query(sql_user_albums_se)
		fmt.Println(res)
		matchedUserAlbumCount = *res[0][0].LongValue
		if matchedUserAlbumCount > 0 {
			return nil, errors.New(fmt.Sprintf("user already owns a path containing this path: %s", rootPath))
		}

		//if err := db.Model(&owner).Association("Albums").Append(&album); err != nil { // INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`,`id`) VALUES ('2022-07-27 03:19:43.834','2022-07-27 03:19:43.834','chengbian',NULL,'/Users/chengbian','1f7d46bead15a13c555e549e19624753',NULL,108) ON DUPLICATE KEY UPDATE `id`=`id`     INSERT INTO `user_albums` (`user_id`,`album_id`) VALUES (18,108) ON DUPLICATE KEY UPDATE `user_id`=`user_id`       UPDATE `users` SET `updated_at`='2022-07-28 10:07:37.749' WHERE `id` = 13
		//	return nil, errors.Wrap(err, "add owner to already existing album")
		//}
		//var albumP string
		//if album.ParentAlbumID == nil {
		//	albumP = "NULL"
		//} else {
		//	albumP = strconv.Itoa(*album.ParentAlbumID)
		//}
		//var albumC string
		//if album.CoverID == nil {
		//	albumC = "NULL"
		//} else {
		//	albumC = strconv.Itoa(*album.ParentAlbumID)
		//}
		//sql_albums_in := "INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`,`id`) VALUES (,'2022-07-27 03:19:43.834','chengbian',NULL,'/Users/chengbian','1f7d46bead15a13c555e549e19624753',NULL,108) ON DUPLICATE KEY UPDATE `id`=`id`"
		//sql_albums_in := fmt.Sprintf("INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`,`id`) VALUES ('%v','%v','%v',%v,'%v','%v',%v,%v) ON DUPLICATE KEY UPDATE `id`=`id`", album.CreatedAt.Format("2006-01-02 15:04:05"), album.UpdatedAt.Format("2006-01-02 15:04:05"), album.Title, albumP, album.Path, album.PathHash, albumC, album.ID)
		////sql_user_albums_in := "INSERT INTO `user_albums` (`user_id`,`album_id`) VALUES (18,108) ON DUPLICATE KEY UPDATE `user_id`=`user_id`"
		sql_user_albums_in := fmt.Sprintf("INSERT INTO `user_albums` (`user_id`,`album_id`) VALUES (%v,%v)", owner.ID, album.ID)
		//sql_users_up := "UPDATE `users` SET `updated_at`=NOW WHERE `id` = 13"
		sql_users_up := fmt.Sprintf("UPDATE `users` SET `updated_at`=NOW WHERE `id` = %v", owner.ID)
		//dataApi.ExecuteSQl(sql_albums_in)
		dataApi.ExecuteSQl(sql_user_albums_in)
		dataApi.ExecuteSQl(sql_users_up)

		return &album, nil
	} else {
		album := models.Album{
			Title:  path.Base(rootPath),
			Path:   rootPath,
			Owners: owners,
		}

		//if err := db.Debug().Create(&album).Error; err != nil { //INSERT INTO `users` (`created_at`,`updated_at`,`username`,`password`,`admin`,`id`) VALUES ('2022-07-29 16:53:26.95','2022-07-29 16:53:26.95','dvaergher',NULL,true,20) ON DUPLICATE KEY UPDATE `id`=`id`      INSERT INTO `user_albums` (`album_id`,`user_id`) VALUES (112,20) ON DUPLICATE KEY UPDATE `album_id`=`album_id`     INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`) VALUES ('2022-07-29 16:53:26.813','2022-07-29 16:53:26.813','suiapp',NULL,'/Users/chengbian/suiapp','15ee3b09f80f4790477703d6ef00f4de',NULL)
		//	return nil, err
		//}
		//sql_users_in := "INSERT INTO `users` (`created_at`,`updated_at`,`username`,`password`,`admin`,`id`) VALUES ('2022-07-29 16:53:26.95','2022-07-29 16:53:26.95','dvaergher',NULL,true,20) ON DUPLICATE KEY UPDATE `id`=`id`"
		//sql_user_albums_in := "INSERT INTO `user_albums` (`album_id`,`user_id`) VALUES (112,20) ON DUPLICATE KEY UPDATE `album_id`=`album_id`"
		//sql_users_in:=fmt.Sprintf("INSERT INTO `users` (`created_at`,`updated_at`,`username`,`password`,`admin`,`id`) VALUES ('%v','%v','%v','%v',%v,%v) ON DUPLICATE KEY UPDATE `id`=`id`",owner.CreatedAt.Format("2006-01-02 15:04:05"),owner.UpdatedAt.Format("2006-01-02 15:04:05"),owner)

		//sql_albums_in := "INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`) VALUES ('2022-07-29 16:53:26.813','2022-07-29 16:53:26.813','suiapp',NULL,'/Users/chengbian/suiapp','15ee3b09f80f4790477703d6ef00f4de',NULL)"
		sql_albums_in := fmt.Sprintf("INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`) VALUES (NOW(),NOW(),'%v',NULL,'%v','%v',NULL) ON DUPLICATE KEY UPDATE `id`=`id`", album.Title, album.Path, models.MD5Hash(album.Path))
		////dataApi.ExecuteSQl(sql_users_in)
		sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE path_hash ='%v'", models.MD5Hash(album.Path))
		dataApi.ExecuteSQl(sql_albums_in)
		res, err = dataApi.Query(sql_albums_se)
		id := DataApi.GetInt(res, 0, 0)
		sql_user_albums_in := fmt.Sprintf("INSERT INTO `user_albums` (`user_id`,`album_id`) VALUES (%v,%v) ", owner.ID, id)
		dataApi.ExecuteSQl(sql_user_albums_in)
		album.ID = id
		fmt.Println("id", id)
		album.PathHash = models.MD5Hash(album.Path)
		return &album, nil
	}
}

var ErrorInvalidRootPath = errors.New("invalid root path")

func ValidRootPath(rootPath string) bool {
	_, err := os.Stat(rootPath)
	if err != nil {
		log.Printf("Warn: invalid root path: '%s'\n%s\n", rootPath, err)
		return false
	}

	return true
}
func ScanAlbum(ctx scanner_task.TaskContext) error {

	newCtx, err := scanner_tasks.Tasks.BeforeScanAlbum(ctx)
	if err != nil {
		return errors.Wrapf(err, "before scan album (%s)", ctx.GetAlbum().Path)
	}
	ctx = newCtx

	// Scan for photos
	albumMedia := findMediaForAlbum(ctx)
	if err != nil {
		return errors.Wrapf(err, "find media for album (%s): %s", ctx.GetAlbum().Path, err)
	}

	changedMedia := make([]*models.Media, 0)
	for i, media := range albumMedia {
		updatedURLs := []*models.MediaURL{}

		mediaData := media_encoding.NewEncodeMediaData(media)

		// define new ctx for scope of for-loop
		ctx, err := scanner_tasks.Tasks.BeforeProcessMedia(ctx, &mediaData)
		if err != nil {
			return err
		}

		//transactionError := ctx.DatabaseTransaction(func(ctx scanner_task.TaskContext) error {
		//	updatedURLs, err = processMedia(ctx, &mediaData)
		//	if err != nil {
		//		return errors.Wrapf(err, "process media (%s)", media.Path)
		//	}
		//
		//	if len(updatedURLs) > 0 {
		//		changedMedia = append(changedMedia, media)
		//	}
		//
		//	return nil
		//})

		updatedURLs, err = processMedia(ctx, &mediaData)
		if err != nil {
			return errors.Wrapf(err, "process media (%s)", media.Path)
		}

		if len(updatedURLs) > 0 {

			changedMedia = append(changedMedia, media)
		}

		//return nil

		//if transactionError != nil {
		//	return errors.Wrap(err, "process media database transaction")
		//}
		if err = scanner_tasks.Tasks.AfterProcessMedia(ctx, &mediaData, updatedURLs, i, len(albumMedia)); err != nil {
			return errors.Wrap(err, "after process media")
		}
	}
	if err := scanner_tasks.Tasks.AfterScanAlbum(ctx, changedMedia, albumMedia); err != nil {
		return errors.Wrap(err, "after scan album")
	}

	return nil
}

func findMediaForAlbum(ctx scanner_task.TaskContext) []*models.Media {

	albumMedia := make([]*models.Media, 0)

	dirContent, err := ioutil.ReadDir(ctx.GetAlbum().Path)
	if err != nil {
		return nil
	}

	for _, item := range dirContent {
		mediaPath := path.Join(ctx.GetAlbum().Path, item.Name())

		isDirSymlink, err := utils.IsDirSymlink(mediaPath)
		if err != nil {
			log.Printf("Cannot detect whether %s is symlink to a directory. Pretending it is not", mediaPath)
			isDirSymlink = false
		}

		if !item.IsDir() && !isDirSymlink && ctx.GetCache().IsPathMedia(mediaPath) {
			skip, err := scanner_tasks.Tasks.MediaFound(ctx, item, mediaPath)
			if err != nil {
				return nil
			}
			if skip {
				continue
			}

			//err = ctx.DatabaseTransaction(func(ctx scanner_task.TaskContext) error {
			//	media, isNewMedia, err := ScanMedia(ctx.GetDB(), mediaPath, ctx.GetAlbum().ID, ctx.GetCache())
			//	if err != nil {
			//		return errors.Wrapf(err, "scanning media error (%s)", mediaPath)
			//	}
			//
			//	if err = scanner_tasks.Tasks.AfterMediaFound(ctx, media, isNewMedia); err != nil {
			//		return err
			//	}
			//
			//	albumMedia = append(albumMedia, media)
			//
			//	return nil
			//})
			{
				id := ctx.GetAlbum().ID
				media, isNewMedia, err := ScanMedia( /*ctx.GetDB(), */ mediaPath, id, ctx.GetCache())
				if err != nil {
					return nil
				}

				scanner_tasks.Tasks.AfterMediaFound(ctx, media, isNewMedia)

				albumMedia = append(albumMedia, media)

				//return nil
			}
			if err != nil {
				scanner_utils.ScannerError("Error scanning media for album (%d): %s\n", ctx.GetAlbum().ID, err)
				continue
			}

			if err != nil {
				scanner_utils.ScannerError("Error scanning media for album (%d): %s\n", ctx.GetAlbum().ID, err)
				continue
			}
		}

	}

	return albumMedia
}

func processMedia(ctx scanner_task.TaskContext, mediaData *media_encoding.EncodeMediaData) ([]*models.MediaURL, error) {

	// Make sure media cache directory exists
	mediaCachePath, err := mediaData.Media.CachePath()
	if err != nil {
		return []*models.MediaURL{}, errors.Wrap(err, "cache directory error")
	}

	return scanner_tasks.Tasks.ProcessMedia(ctx, mediaData, mediaCachePath)
}
