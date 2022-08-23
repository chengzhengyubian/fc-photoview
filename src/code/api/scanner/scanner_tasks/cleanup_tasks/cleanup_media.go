package cleanup_tasks

import (
	"encoding/json"
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/scanner_utils"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
)

// CleanupMedia removes media entries from the database that are no longer present on the filesystem
func CleanupMedia( /*db *gorm.DB, */ albumId int, albumMedia []*models.Media) []error {
	albumMediaIds := make([]int, len(albumMedia))
	for i, media := range albumMedia {
		albumMediaIds[i] = media.ID
	}
	fmt.Println(albumId, albumMediaIds)
	var sql_media_se string
	// Will get from database
	var mediaList []models.Media

	albumMediaId, err := json.Marshal(albumMediaIds)
	fmt.Println(err)
	albumMediaids := strings.Trim(string(albumMediaId), "[]")
	// Select media from database that was not found on hard disk
	if len(albumMedia) > 0 {
		sql_media_se = fmt.Sprintf("SELECT * FROM `media` WHERE album_id = %v AND NOT id IN (%v)", albumId, albumMediaids)
	} else {
		sql_media_se = fmt.Sprintf("SELECT * FROM `media` WHERE album_id = %v", albumId)
	}
	dataApi, err := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var Media models.Media
		Media.ID = dataapi.GetInt(res, i, 0)
		Media.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		Media.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		Media.Title = *res[i][3].StringValue
		Media.Path = *res[i][4].StringValue
		Media.PathHash = *res[i][5].StringValue
		Media.AlbumID = int(*res[i][6].LongValue)
		Media.ExifID = dataapi.GetIntP(res, i, 7)
		Media.DateShot = time.Unix(*res[i][8].LongValue/1000, 0)
		if *res[0][9].StringValue == "photo" {
			Media.Type = models.MediaTypePhoto
		} else {
			Media.Type = models.MediaTypeVideo
		}
		Media.VideoMetadataID = dataapi.GetIntP(res, i, 10)
		Media.SideCarPath = dataapi.GetStringP(res, i, 11)
		Media.SideCarHash = dataapi.GetStringP(res, i, 12)
		Media.Blurhash = dataapi.GetStringP(res, i, 13)
		mediaList = append(mediaList, Media)
	}
	deleteErrors := make([]error, 0)

	mediaIDs := make([]int, 0)
	for _, media := range mediaList {

		mediaIDs = append(mediaIDs, media.ID)
		cachePath := path.Join(utils.MediaCachePath(), strconv.Itoa(int(albumId)), strconv.Itoa(int(media.ID)))
		err := os.RemoveAll(cachePath)
		if err != nil {
			deleteErrors = append(deleteErrors, errors.Wrapf(err, "delete unused cache folder (%s)", cachePath))
		}

	}

	if len(mediaIDs) > 0 {
		mediaID, err := json.Marshal(mediaIDs)
		fmt.Println(err)
		mediaids := strings.Trim(string(mediaID), "[]")
		sql_media_de := fmt.Sprintf("DELETE FROM `media` WHERE id IN (%v)", mediaids)
		dataApi.ExecuteSQl(sql_media_de)
		fmt.Println("删除media", mediaids)
	}

	return deleteErrors
}

// DeleteOldUserAlbums finds and deletes old albums in the database and cache that does not exist on the filesystem anymore.
func DeleteOldUserAlbums(scannedAlbums []*models.Album, user *models.User) []error {
	if len(scannedAlbums) == 0 {
		return nil
	}

	scannedAlbumIDs := make([]interface{}, len(scannedAlbums))
	for i, album := range scannedAlbums {
		scannedAlbumIDs[i] = album.ID
	}

	// Old albums to be deleted
	var deleteAlbums []models.Album
	scannedAlbumID, err := json.Marshal(scannedAlbumIDs)
	fmt.Println(err)
	scannedAlbumids := strings.Trim(string(scannedAlbumID), "[]")
	sql_albums_se := fmt.Sprintf("SELECT albums.* FROM `user_albums` JOIN albums ON user_albums.album_id = albums.id WHERE user_id = %v AND album_id NOT IN (%v)", user.ID, scannedAlbumids)
	dataApi, _ := dataapi.NewDataApiClient()
	res, _ := dataApi.Query(sql_albums_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var album models.Album
		album.ID = dataapi.GetInt(res, i, 0)
		album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		album.Title = dataapi.GetString(res, i, 3)
		album.ParentAlbumID = dataapi.GetIntP(res, i, 4)
		album.Path = dataapi.GetString(res, i, 5)
		album.PathHash = dataapi.GetString(res, i, 6)
		album.CoverID = dataapi.GetIntP(res, i, 7)
		deleteAlbums = append(deleteAlbums, album)
	}
	if len(deleteAlbums) == 0 {
		return []error{}
	}

	deleteErrors := make([]error, 0)

	// Delete old albums from cache
	deleteAlbumIDs := make([]int, len(deleteAlbums))
	for i, album := range deleteAlbums {
		deleteAlbumIDs[i] = album.ID
		cachePath := path.Join(utils.MediaCachePath(), strconv.Itoa(int(album.ID)))
		err := os.RemoveAll(cachePath)
		if err != nil {
			deleteErrors = append(deleteErrors, errors.Wrapf(err, "delete unused cache folder (%s)", cachePath))
		}
	}
	//删除user_albums
	deleteAlbumID, err := json.Marshal(deleteAlbumIDs)
	fmt.Println(err)
	deleteAlbumids := strings.Trim(string(deleteAlbumID), "[]")
	sql_User_albums_de := fmt.Sprintf("DELETE FROM `user_albums` WHERE album_id IN (%v)", deleteAlbumids)
	dataApi.ExecuteSQl(sql_User_albums_de)
	//删除albums
	sql_albums_de := fmt.Sprintf("DELETE FROM `albums` WHERE id IN (%v)", deleteAlbumids)
	dataApi.ExecuteSQl(sql_albums_de)

	if err != nil {
		scanner_utils.ScannerError("Could not delete old albums from database:\n%s\n", err)
		deleteErrors = append(deleteErrors, err)
	}

	return deleteErrors
}
