package actions

//修改完
import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
)

//修改完，测试成功
func MyTimeline(user *models.User, paginate *models.Pagination, onlyFavorites *bool, fromDate *time.Time) ([]*models.Media, error) {

	//query := db.
	//	Joins("JOIN albums ON media.album_id = albums.id").
	//	Where("albums.id IN (?)", db.Table("user_albums").Select("user_albums.album_id").Where("user_id = ?", user.ID))
	//
	//switch drivers.GetDatabaseDriverType(db) {
	//case drivers.POSTGRES:
	//	query = query.
	//		Order("DATE_TRUNC('year', date_shot) DESC").
	//		Order("DATE_TRUNC('month', date_shot) DESC").
	//		Order("DATE_TRUNC('day', date_shot) DESC").
	//		Order("albums.title ASC").
	//		Order("media.date_shot DESC")
	//case drivers.SQLITE:
	//	query = query.
	//		Order("strftime('%j', media.date_shot) DESC"). // convert to day of year 001-366
	//		Order("albums.title ASC").
	//		Order("TIME(media.date_shot) DESC")
	//default:
	//	query = query.
	//		Order("YEAR(media.date_shot) DESC").
	//		Order("MONTH(media.date_shot) DESC").
	//		Order("DAY(media.date_shot) DESC").
	//		Order("albums.title ASC").
	//		Order("TIME(media.date_shot) DESC")
	//}

	//从某天开始，修改中，未测试
	var limit int
	var offset int
	limit = *paginate.Limit
	offset = *paginate.Offset
	var sql_media_se string
	if fromDate != nil {
		if onlyFavorites != nil && *onlyFavorites {
			sql_media_se = fmt.Sprintf("SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = %v) AND media.date_shot < '%v' AND media.id IN (SELECT user_media_data.media_id FROM `user_media_data` WHERE user_media_data.user_id = %v AND user_media_data.favorite) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT %v OFFSET %v", user.ID, fromDate, user.ID, limit, offset)
			//query = query.Where("media.id IN (?)", db.Table("user_media_data").Select("user_media_data.media_id").Where("user_media_data.user_id = ?", user.ID).Where("user_media_data.favorite"))
		} else {
			sql_media_se = fmt.Sprintf("SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = %v) AND media.date_shot < '%v' ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT %v OFFSET %v", user.ID, fromDate, limit, offset)
		}
	} else {
		if onlyFavorites != nil && *onlyFavorites {
			sql_media_se = fmt.Sprintf("SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = %v) AND media.id IN (SELECT user_media_data.media_id FROM `user_media_data` WHERE user_media_data.user_id = %v AND user_media_data.favorite) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT %v OFFSET %v", user.ID, user.ID, limit, offset)
		} else {
			sql_media_se = fmt.Sprintf("SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = %v) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT %v OFFSET %v", user.ID, limit, offset)
		}
	}
	//query = models.FormatSQL(query, nil, paginate)
	//SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = 2) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT 200 OFFSET 3
	var media []*models.Media //只喜欢
	//if err := query.Find(&media).Error; err != nil { //SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = 2) AND media.id IN (SELECT user_media_data.media_id FROM `user_media_data` WHERE user_media_data.user_id = 2 AND user_media_data.favorite) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT 200 OFFSET 1
	//	return nil, err
	//}
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var Media models.Media
		Media.ID = DataApi.GetInt(res, i, 0)
		Media.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		Media.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		Media.Title = *res[i][3].StringValue
		Media.Path = *res[i][4].StringValue
		Media.PathHash = *res[i][5].StringValue
		Media.AlbumID = int(*res[i][6].LongValue)
		Media.ExifID = DataApi.GetIntP(res, i, 7)
		Media.DateShot = time.Unix(*res[i][8].LongValue/1000, 0)
		if *res[0][9].StringValue == "photo" {
			Media.Type = models.MediaTypePhoto
		} else {
			Media.Type = models.MediaTypeVideo
		}
		Media.VideoMetadataID = DataApi.GetIntP(res, i, 10)
		Media.SideCarPath = DataApi.GetStringP(res, i, 11)
		Media.SideCarHash = DataApi.GetStringP(res, i, 12)
		Media.Blurhash = DataApi.GetStringP(res, i, 13)
		media = append(media, &Media)
	}
	if err != nil {
		return nil, err
	}
	return media, nil

	//只显示喜爱
	//SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = 2) AND media.date_shot < '2022-01-01 00:00:00' AND media.id IN (SELECT user_media_data.media_id FROM `user_media_data` WHERE user_media_data.user_id = 2 AND user_media_data.favorite) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT 200 OFFSET 1
	//从今天开始
	//SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash` FROM `media` JOIN albums ON media.album_id = albums.id WHERE albums.id IN (SELECT user_albums.album_id FROM `user_albums` WHERE user_id = 2) ORDER BY YEAR(media.date_shot) DESC,MONTH(media.date_shot) DESC,DAY(media.date_shot) DESC,albums.title ASC,TIME(media.date_shot) DESC LIMIT 200 OFFSET 3
}
