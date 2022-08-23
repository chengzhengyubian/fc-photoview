package actions

import (
	"github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/models"
	"strconv"
	"time"
)

func MyMedia(user *models.User, order *models.Ordering, paginate *models.Pagination) ([]*models.Media, error) {
	if err := user.FillAlbums(); err != nil {
		return nil, err
	}
	var media []*models.Media
	sql_media_se := "SELECT * FROM `media` WHERE media.album_id IN (SELECT user_albums.album_id FROM user_albums WHERE user_albums.user_id = " + strconv.Itoa(user.ID) + ") ORDER BY " + *order.OrderBy + " LIMIT " + strconv.Itoa(*paginate.Limit)
	dataApi, _ := dataapi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	if num == 0 {
		return nil, err
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
	return media, nil
}
