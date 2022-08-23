package scanner

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
		sql_user_albums_se := fmt.Sprintf("SELECT count(*) FROM `user_albums` WHERE user_id =%v and album_id = %v", owner.ID, album.ID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err = dataApi.Query(sql_user_albums_se)
		fmt.Println(res)
		matchedUserAlbumCount = *res[0][0].LongValue
		if matchedUserAlbumCount > 0 {
			return nil, errors.New(fmt.Sprintf("user already owns a path containing this path: %s", rootPath))
		}
		sql_user_albums_in := fmt.Sprintf("INSERT INTO `user_albums` (`user_id`,`album_id`) VALUES (%v,%v)", owner.ID, album.ID)
		sql_users_up := fmt.Sprintf("UPDATE `users` SET `updated_at`=NOW WHERE `id` = %v", owner.ID)
		dataApi.ExecuteSQl(sql_user_albums_in)
		dataApi.ExecuteSQl(sql_users_up)

		return &album, nil
	} else {
		album := models.Album{
			Title:  path.Base(rootPath),
			Path:   rootPath,
			Owners: owners,
		}
		sql_albums_in := fmt.Sprintf("INSERT INTO `albums` (`created_at`,`updated_at`,`title`,`parent_album_id`,`path`,`path_hash`,`cover_id`) VALUES (NOW(),NOW(),'%v',NULL,'%v','%v',NULL) ON DUPLICATE KEY UPDATE `id`=`id`", album.Title, album.Path, models.MD5Hash(album.Path))
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
		updatedURLs, err = processMedia(ctx, &mediaData)
		if err != nil {
			return errors.Wrapf(err, "process media (%s)", media.Path)
		}

		if len(updatedURLs) > 0 {

			changedMedia = append(changedMedia, media)
		}
		if err = scanner_tasks.Tasks.AfterProcessMedia(ctx, &mediaData, updatedURLs, i, len(albumMedia)); err != nil {
			return errors.Wrap(err, "after process media")
		}
	}
	DataApi.Stresstest()
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
