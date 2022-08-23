package actions

import (
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"

	//"github.com/photoview/photoview/api/database/drivers"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
	//"gorm.io/gorm"
	//"gorm.io/gorm/clause"
)

func Search(query string, userID int, _limitMedia *int, _limitAlbums *int) (*models.SearchResult, error) {
	limitMedia := 10
	limitAlbums := 10

	if _limitMedia != nil {
		limitMedia = *_limitMedia
	}

	if _limitAlbums != nil {
		limitAlbums = *_limitAlbums
	}

	wildQuery := "%" + strings.ToLower(query) + "%"

	var media []*models.Media
	var sql_user_albums_select string

	sql_user_albums_select = fmt.Sprintf("select * from user_albums where user_id =%v", userID)
	sql_media_se := fmt.Sprintf("select media.* from media left join albums on media.album_id=albums.id where EXISTS (%v) and LOWER(media.title) LIKE '%v' OR LOWER(media.path) LIKE '%v' order by  (CASE WHEN LOWER(media.title) LIKE '%v' THEN 2 WHEN LOWER(media.path) LIKE '%v' THEN 1 END) DESC limit %v", sql_user_albums_select, wildQuery, wildQuery, wildQuery, wildQuery, limitMedia)
	dataApi, _ := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	if num == 0 {
		return nil, errors.Wrapf(err, "searching media")
	}
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
		media = append(media, &Media)
	}
	var albums []*models.Album
	sql_albums_se := fmt.Sprintf("select * from albums where EXISTS (select * from user_albums where user_id=%v and user_albums.album_id=albums.id) and albums.title LIKE '%v' OR albums.path LIKE '%v' order by (CASE WHEN albums.title LIKE '%v' THEN 2 WHEN albums.path LIKE '%v' THEN 1 END) DESC limit %v", userID, wildQuery, wildQuery, wildQuery, wildQuery, limitAlbums)
	res, err = dataApi.Query(sql_albums_se)
	num = len(res)
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
		albums = append(albums, &album)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "searching albums")
	}

	result := models.SearchResult{
		Query:  query,
		Media:  media,
		Albums: albums,
	}

	return &result, nil
}
