package scanner

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

func ScanMedia(mediaPath string, albumId int, cache *scanner_cache.AlbumScannerCache) (*models.Media, bool, error) {
	mediaName := path.Base(mediaPath)

	{ // Check if media already exists
		var media []*models.Media
		sql_media_se := fmt.Sprintf("select * from media where path_hash='%v'", models.MD5Hash(mediaPath))
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_se)
		if err != nil {
			return nil, false, errors.Wrap(err, "scan media fetch from database")
		}
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
	var Type string
	if media.Type == models.MediaTypePhoto {
		Type = "photo"
	} else {
		Type = "video"
	}
	sql_media_in := fmt.Sprintf("insert into media (created_at, updated_at,title,path, path_hash,album_id, date_shot,type) values(NOW(),NOW(),'%v','%v','%v',%v,'%v','%v')", media.Title, media.Path, models.MD5Hash(media.Path), media.AlbumID, timestr, Type)
	sql_media_se := fmt.Sprintf("select id from media where path_hash='%v'", models.MD5Hash(media.Path))
	dataApi, _ := DataApi.NewDataApiClient()
	//体现serverless特性
	//sql_serverless_test := "select benchmark(37000000 ,crc32('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef'))"
	//dataApi.Query(sql_serverless_test)
	dataApi.ExecuteSQl(sql_media_in)
	res, err := dataApi.Query(sql_media_se)
	for len(res) == 0 {
		dataApi.ExecuteSQl(sql_media_in)
		res, err = dataApi.Query(sql_media_se)
	}
	media.ID = DataApi.GetInt(res, 0, 0)
	return &media, true, nil
}

// ProcessSingleMedia processes a single media, might be used to reprocess media with corrupted cache
// Function waits for processing to finish before returning.
func ProcessSingleMedia( /*db *gorm.DB, */ media *models.Media) error {
	album_cache := scanner_cache.MakeAlbumCache()

	var album models.Album
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
