package scanner

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"image"
	"log"
	"os"
	"time"

	"github.com/buckket/go-blurhash"
	"github.com/photoview/photoview/api/graphql/models"
)

// GenerateBlurhashes queries the database for media that are missing a blurhash and computes one for them.
// This function blocks until all hashes have been computed
func GenerateBlurhashes( /*db *gorm.DB*/ ) error {
	var results []*models.Media

	processErrors := make([]error, 0)
	sql_media_se := fmt.Sprintf("SELECT `media`.`id`,`media`.`created_at`,`media`.`updated_at`,`media`.`title`,`media`.`path`,`media`.`path_hash`,`media`.`album_id`,`media`.`exif_id`,`media`.`date_shot`,`media`.`type`,`media`.`video_metadata_id`,`media`.`side_car_path`,`media`.`side_car_hash`,`media`.`blurhash`,media_urls.id,media_urls.created_at,media_urls.updated_at,media_urls.media_id,media_urls.media_name,media_urls.width,media_urls.height,media_urls.purpose,media_urls.content_type,media_urls.file_size FROM `media` INNER JOIN media_urls ON media.id = media_urls.media_id WHERE blurhash IS NULL AND (media_urls.purpose = 'thumbnail' OR media_urls.purpose = 'video-thumbnail') ORDER BY `media`.`id` LIMIT 50")
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var Media models.Media
		var MediaURL models.MediaURL
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
		MediaURL.ID = DataApi.GetInt(res, i, 14)
		MediaURL.CreatedAt = time.Unix(*res[i][15].LongValue/1000, 0)
		MediaURL.UpdatedAt = time.Unix(*res[i][16].LongValue/1000, 0)
		MediaURL.MediaID = DataApi.GetInt(res, i, 17)
		MediaURL.MediaName = DataApi.GetString(res, i, 18)
		MediaURL.Width = DataApi.GetInt(res, i, 19)
		MediaURL.Height = DataApi.GetInt(res, i, 20)
		switch DataApi.GetString(res, i, 21) {
		case "thumbnail":
			MediaURL.Purpose = models.PhotoThumbnail
		case "high-res":
			MediaURL.Purpose = models.PhotoHighRes
		case "original":
			MediaURL.Purpose = models.MediaOriginal
		case "video-web":
			MediaURL.Purpose = models.VideoWeb
		case "video-thumbnail":
			MediaURL.Purpose = models.VideoThumbnail
		}
		MediaURL.ContentType = DataApi.GetString(res, i, 22)
		MediaURL.FileSize = DataApi.GetLong(res, i, 23)
		Media.MediaURL = append(Media.MediaURL, MediaURL)
		results = append(results, &Media)
	}
	log.Printf("generating %d blurhashes", len(results))

	hashes := make([]*string, len(results))

	for i, row := range results {

		thumbnail, err := row.GetThumbnail()
		if err != nil {
			log.Printf("failed to get thumbnail for media to generate blurhash (%d): %v", row.ID, err)
			processErrors = append(processErrors, err)
			continue
		}

		hashStr, err := GenerateBlurhashFromThumbnail(thumbnail)
		if err != nil {
			log.Printf("failed to generate blurhash for media (%d): %v", row.ID, err)
			processErrors = append(processErrors, err)
			continue
		}

		hashes[i] = &hashStr
		results[i].Blurhash = &hashStr
		fmt.Println(hashStr)
		sql_media_up := fmt.Sprintf("update media set updated_at=NOW(),blurhash='%v' where id=%v", hashStr, row.ID)
		dataApi.ExecuteSQl(sql_media_up)
	}

	if err != nil {
		return err
	}

	if len(processErrors) == 0 {
		return nil
	} else {
		return fmt.Errorf("failed to generate %d blurhashes", len(processErrors))
	}
}

// GenerateBlurhashFromThumbnail generates a blurhash for a single media and stores it in the database
func GenerateBlurhashFromThumbnail(thumbnail *models.MediaURL) (string, error) {
	thumbnail_path, err := thumbnail.CachedPath()
	if err != nil {
		return "", err
	}

	imageFile, err := os.Open(thumbnail_path)
	if err != nil {
		return "", err
	}

	imageData, _, err := image.Decode(imageFile)
	if err != nil {
		return "", err
	}

	hashStr, err := blurhash.Encode(4, 3, imageData)
	if err != nil {
		return "", err
	}

	return hashStr, nil
}
