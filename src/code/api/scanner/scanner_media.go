package scanner

/*修改完*/
import (
	"context"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/scanner_cache"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/scanner/scanner_tasks"
	"github.com/pkg/errors"
	//"gorm.io/gorm"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

//

func ScanMedia( /*tx *gorm.DB, */ mediaPath string, albumId int, cache *scanner_cache.AlbumScannerCache) (*models.Media, bool, error) {
	mediaName := path.Base(mediaPath)

	{ // Check if media already exists
		var media []*models.Media
		//result := tx.Debug().Where("path_hash = ?", models.MD5Hash(mediaPath)).Find(&media)
		//sql_media_se := "select * from media where path_hash =\"" + models.MD5Hash(mediaPath) + "\""
		sql_media_se := fmt.Sprintf("select * from media where path_hash='%v'", models.MD5Hash(mediaPath))
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_se)
		if err != nil {
			return nil, false, errors.Wrap(err, "scan media fetch from database")
		}
		//if len(res) == 0 {
		//	return nil, false, errors.Wrap(err, "scan media fetch from database")
		//}
		//if result.Error != nil {
		//	return nil, false, errors.Wrap(result.Error, "scan media fetch from database")
		//}
		fmt.Print(res, "media result")
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
		if len(res) > 0 {
			// log.Printf("Media already scanned: %s\n", mediaPath)
			return media[0], false, nil
		}
		//if result.RowsAffected > 0 {
		//	// log.Printf("Media already scanned: %s\n", mediaPath)
		//	return media[0], false, nil
		//}
	}
	log.Printf("Scanning media: %s\n", mediaPath)

	mediaType, err := cache.GetMediaType(mediaPath)
	if err != nil {
		return nil, false, errors.Wrap(err, "could determine if media was photo or video")
	}

	var mediaTypeText models.MediaType

	if mediaType.IsVideo() {
		mediaTypeText = models.MediaTypeVideo
	} else {
		mediaTypeText = models.MediaTypePhoto
	}

	stat, err := os.Stat(mediaPath)
	if err != nil {
		return nil, false, err
	}

	media := models.Media{
		Title:    mediaName,
		Path:     mediaPath,
		AlbumID:  albumId,
		Type:     mediaTypeText,
		DateShot: stat.ModTime(),
	}
	timestr := media.DateShot.Format("2006-01-02 15:04:05")
	/*if err := tx.Debug().Create(&media).Error; err != nil { //INSERT INTO `media` (`created_at`,`updated_at`,`title`,`path`,`path_hash`,`album_id`,`exif_id`,`date_shot`,`type`,`video_metadata_id`,`side_car_path`,`side_car_hash`,`blurhash`) VALUES ('2022-08-04 11:18:51.083','2022-08-04 11:18:51.083','自我介绍.png','/Users/chengbian/suiapp/photo/photoper/自我介绍.png','f9b6fd9b47d427178a9d22198dac7c9c',148,NULL,'2022-06-30 17:50:25.916','photo',NULL,NULL,NULL,NULL
		return nil, false, errors.Wrap(err, "could not insert media into database")
	}*/
	var Type string
	if media.Type == models.MediaTypePhoto {
		Type = "photo"
	} else {
		Type = "video"
	}
	//sql_media_in := "insert into media (created_at, updated_at,title,path, path_hash,album_id, date_shot,type) values(NOW(),NOW(),\"" + media.Title + "\",\"" + media.Path + "\",\"" + models.MD5Hash(media.Path) + "\"," + strconv.Itoa(media.AlbumID) + ",\"" + timestr + "\",\"" + Type + "\")"
	sql_media_in := fmt.Sprintf("insert into media (created_at, updated_at,title,path, path_hash,album_id, date_shot,type) values(NOW(),NOW(),'%v','%v','%v',%v,'%v','%v')", media.Title, media.Path, models.MD5Hash(media.Path), media.AlbumID, timestr, Type)
	//sql_media_urls_in := fmt.Sprintf("insert into media_urls(created_at, updated_at,media_id,media_name,)")
	sql_media_se := fmt.Sprintf("select id from media where path_hash='%v'", models.MD5Hash(media.Path))
	dataApi, _ := DataApi.NewDataApiClient()

	dataApi.ExecuteSQl(sql_media_in)
	res, err := dataApi.Query(sql_media_se)
	for len(res) == 0 {
		dataApi.ExecuteSQl(sql_media_in)
		res, err = dataApi.Query(sql_media_se)
	}
	media.ID = DataApi.GetInt(res, 0, 0)
	return &media, true, nil
}

//这里有一个函数还没改，已经改了大部分

// ProcessSingleMedia processes a single media, might be used to reprocess media with corrupted cache
// Function waits for processing to finish before returning.
func ProcessSingleMedia( /*db *gorm.DB, */ media *models.Media) error {
	album_cache := scanner_cache.MakeAlbumCache()

	var album models.Album
	//if err := db.Model(media).Association("Album").Find(&album); err != nil { //SELECT * FROM `albums` WHERE `albums`.`id` = 1
	//
	//	return err
	//}
	sql_albums_se := "SELECT * FROM `albums` WHERE `albums`.`id` =" + strconv.Itoa(media.AlbumID)
	dataAPi, _ := DataApi.NewDataApiClient()
	res, err := dataAPi.Query(sql_albums_se)
	if len(res) == 0 {
		return err
	}
	album.ID = int(*res[0][0].LongValue)
	album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	album.Title = *res[0][3].StringValue
	album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
	album.Path = *res[0][5].StringValue
	album.PathHash = *res[0][6].StringValue
	album.CoverID = DataApi.GetIntP(res, 0, 7)
	media_data := media_encoding.NewEncodeMediaData(media)

	task_context := scanner_task.NewTaskContext(context.Background() /*db,*/, &album, album_cache) //这里注意一下，这里还没改
	new_ctx, err := scanner_tasks.Tasks.BeforeProcessMedia(task_context, &media_data)
	if err != nil {
		return err
	}

	mediaCachePath, err := media.CachePath() //
	if err != nil {
		return err
	}

	updated_urls, err := scanner_tasks.Tasks.ProcessMedia(new_ctx, &media_data, mediaCachePath)
	if err != nil {
		return err
	}

	err = scanner_tasks.Tasks.AfterProcessMedia(new_ctx, &media_data, updated_urls, 0, 1)
	if err != nil {
		return err
	}

	return nil
}
