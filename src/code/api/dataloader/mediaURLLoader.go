package dataloader

/*修改完*/
import (
	"encoding/json"
	"fmt"
	"github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_type"
	"github.com/pkg/errors"
	//"gorm.io/gorm"
)

//修改完，暂时没问题
func makeMediaURLLoader( /*db *gorm.DB,*/ filter func(sql string) string) func(keys []int) ([]*models.MediaURL, []error) {
	return func(mediaIDs []int) ([]*models.MediaURL, []error) {

		var urls []*models.MediaURL
		//query := db.Where("media_id IN (?)", mediaIDs)
		mediaID, _ := json.Marshal(mediaIDs)
		mediaids := strings.Trim(string(mediaID), "[]")
		sql_media_urls_se := fmt.Sprintf("select * from media_urls where media_id IN (%v)", mediaids)
		//filter(query)
		filter(sql_media_urls_se)
		//if err := query.Find(&urls).Error; err != nil { //SELECT * FROM `media_urls` WHERE media_id IN (2) AND (purpose = 'high-res' OR (purpose = 'original' AND content_type IN ('image/jpeg','image/png','image/webp','image/bmp')))
		//	return nil, []error{errors.Wrap(err, "media url loader database query")}
		//}
		dataApi, _ := dataapi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_urls_se)
		if err != nil {
			return nil, []error{errors.Wrap(err, "media url loader database query")}
		}
		num := len(res)
		for i := 0; i < num; i++ {
			var url models.MediaURL
			url.ID = dataapi.GetInt(res, i, 0)
			url.CreatedAt = time.Unix(dataapi.GetLong(res, i, 1)/1000, 0)
			url.UpdatedAt = time.Unix(dataapi.GetLong(res, i, 2)/1000, 0)
			url.MediaID = dataapi.GetInt(res, i, 3)
			url.MediaName = dataapi.GetString(res, i, 4)
			url.Width = dataapi.GetInt(res, i, 5)
			url.Height = dataapi.GetInt(res, i, 6)
			switch dataapi.GetString(res, i, 7) {
			case "thumbnail":
				url.Purpose = models.PhotoThumbnail
			case "high-res":
				url.Purpose = models.PhotoHighRes
			case "original":
				url.Purpose = models.MediaOriginal
			case "video-web":
				url.Purpose = models.VideoWeb
			case "video-thumbnail":
				url.Purpose = models.VideoThumbnail
			}
			url.ContentType = dataapi.GetString(res, i, 8)
			url.FileSize = dataapi.GetLong(res, i, 9)
			urls = append(urls, &url)
		}
		//sql_media_urls_se:="SELECT * FROM `media_urls` WHERE media_id IN (2) AND (purpose = 'high-res' OR (purpose = 'original' AND content_type IN ('image/jpeg','image/png','image/webp','image/bmp')))"
		resultMap := make(map[int]*models.MediaURL, len(mediaIDs))
		for _, url := range urls {
			resultMap[url.MediaID] = url
		}

		result := make([]*models.MediaURL, len(mediaIDs))
		for i, mediaID := range mediaIDs {
			mediaURL, found := resultMap[mediaID]
			if found {
				result[i] = mediaURL
			} else {
				result[i] = nil
			}
		}

		return result, nil
	}
}

func NewThumbnailMediaURLLoader( /*db *gorm.DB*/ ) *MediaURLLoader {
	return &MediaURLLoader{
		maxBatch: 100,
		wait:     5 * time.Millisecond,
		fetch: makeMediaURLLoader( /*db, */ func(sql string) string {
			//return query.Where("purpose = ? OR purpose = ?", models.PhotoThumbnail, models.VideoThumbnail)

			return sql + fmt.Sprintf(" and (purpose = '%v' OR purpose ='%v'", models.PhotoThumbnail, models.VideoThumbnail)
		}),
	}
}

func NewHighresMediaURLLoader( /*db *gorm.DB*/ ) *MediaURLLoader {
	return &MediaURLLoader{
		maxBatch: 100,
		wait:     5 * time.Millisecond,
		fetch: makeMediaURLLoader( /*db, */ func(sql string) string {
			//return query.Where("purpose = ? OR (purpose = ? AND content_type IN ?)", models.PhotoHighRes, models.MediaOriginal, media_type.WebMimetypes)
			return sql + fmt.Sprintf(" and (purpose ='%v' OR (purpose ='%v' AND content_type IN %v)", models.PhotoHighRes, models.MediaOriginal, media_type.WebMimetypes)
		}),
	}
}

func NewVideoWebMediaURLLoader( /*db *gorm.DB*/ ) *MediaURLLoader {
	return &MediaURLLoader{
		maxBatch: 100,
		wait:     5 * time.Millisecond,
		fetch: makeMediaURLLoader( /*db, */ func(sql string) string {
			//return query.Where("purpose = ? OR purpose = ?", models.VideoWeb, models.MediaOriginal)
			return sql + fmt.Sprintf(" and (purpose ='%v' OR (purpose ='%v'", models.VideoWeb, models.MediaOriginal)
		}),
	}
}
