package dataloader

/*修改完*/
import (
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
)

func NewUserFavoriteLoader( /*db *gorm.DB*/ ) *UserFavoritesLoader {
	return &UserFavoritesLoader{
		maxBatch: 100,
		wait:     5 * time.Millisecond,
		fetch: func(keys []*models.UserMediaData) ([]bool, []error) {

			userIDMap := make(map[int]struct{}, len(keys))
			mediaIDMap := make(map[int]struct{}, len(keys))
			for _, key := range keys {
				userIDMap[key.UserID] = struct{}{}
				mediaIDMap[key.MediaID] = struct{}{}
			}

			uniqueUserIDs := make([]int, len(userIDMap))
			uniqueMediaIDs := make([]int, len(mediaIDMap))

			count := 0
			for id := range userIDMap {
				uniqueUserIDs[count] = id
				count++
			}

			count = 0
			for id := range mediaIDMap {
				uniqueMediaIDs[count] = id
				count++
			}
			// SELECT * FROM `user_media_data` WHERE user_id IN (2) AND media_id IN (1) AND favorite = TRUE
			var userMediaFavorites []*models.UserMediaData
			//err := db.Where("user_id IN (?)", uniqueUserIDs).Where("media_id IN (?)", uniqueMediaIDs).Where("favorite = TRUE").Find(&userMediaFavorites).Error
			//if err != nil {
			//	return nil, []error{err}
			//}

			//转为字符串
			uniqueUserID, err := json.Marshal(uniqueUserIDs)
			if err != nil {
				fmt.Print(err)
			}
			uniqueUserids := strings.Trim(string(uniqueUserID), "[]")

			uniqueMediaID, err := json.Marshal(uniqueMediaIDs)
			if err != nil {
				fmt.Print(err)
			}
			uniqueMediaids := strings.Trim(string(uniqueMediaID), "[]")
			sql_user_media_data_se := "SELECT * FROM `user_media_data` WHERE user_id IN (" + uniqueUserids + ") AND media_id IN (" + uniqueMediaids + ") AND favorite = TRUE"
			dataApi, _ := DataApi.NewDataApiClient()
			res, err := dataApi.Query(sql_user_media_data_se)
			num := len(res)
			for i := 0; i < num; i++ {
				var userMediaData models.UserMediaData
				userMediaData.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 0)/1000, 0)
				userMediaData.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
				userMediaData.UserID = DataApi.GetInt(res, 0, 2)
				userMediaData.MediaID = DataApi.GetInt(res, 0, 3)
				userMediaData.Favorite = DataApi.GetBoolean(res, 0, 4)
				userMediaFavorites = append(userMediaFavorites, &userMediaData)
			}
			result := make([]bool, len(keys))
			for i, key := range keys {
				favorite := false
				for _, fav := range userMediaFavorites {
					if fav.UserID == key.UserID && fav.MediaID == key.MediaID {
						favorite = true
						break
					}
				}
				result[i] = favorite
			}

			return result, nil
		},
	}
}
