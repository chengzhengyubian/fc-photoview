package processing_tasks

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/media_encoding/media_utils"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
)

// Higher order function used to check if MediaURL for a given MediaPurpose exists
func makePhotoURLChecker( /*tx *gorm.DB, */ mediaID int) func(purpose models.MediaPurpose) (*models.MediaURL, error) {
	return func(purpose models.MediaPurpose) (*models.MediaURL, error) {
		var mediaURL []*models.MediaURL
		var pur string
		switch purpose {
		case models.PhotoThumbnail:
			pur = "thumbnail"
		case models.VideoWeb:
			pur = "video-web"
		case models.VideoThumbnail:
			pur = "video-thumbnail"
		case models.MediaOriginal:
			pur = "original"
		case models.PhotoHighRes:
			pur = "high-res"
		}
		sql_media_urls_se := "SELECT * FROM `media_urls` WHERE purpose =\"" + pur + "\"AND media_id =" + strconv.Itoa(mediaID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_media_urls_se)
		num := len(res)
		if num == 0 {
			return nil, err
		}
		for i := 0; i < num; i++ {
			var MediaURL models.MediaURL
			MediaURL.ID = DataApi.GetInt(res, i, 0)
			MediaURL.CreatedAt = time.Unix(DataApi.GetLong(res, i, 1)/1000, 0)
			MediaURL.UpdatedAt = time.Unix(DataApi.GetLong(res, i, 2)/1000, 0)
			MediaURL.MediaID = DataApi.GetInt(res, i, 3)
			MediaURL.MediaName = DataApi.GetString(res, i, 4)
			MediaURL.Width = DataApi.GetInt(res, i, 5)
			MediaURL.Height = DataApi.GetInt(res, i, 6)
			switch DataApi.GetString(res, i, 7) {
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
			MediaURL.ContentType = DataApi.GetString(res, i, 8)
			MediaURL.FileSize = DataApi.GetLong(res, i, 9)
			mediaURL = append(mediaURL, &MediaURL)
		}
		if num > 0 {
			return mediaURL[0], nil
		}

		return nil, nil
	}
}

func generateUniqueMediaNamePrefixed(prefix string, mediaPath string, extension string) string {
	mediaName := fmt.Sprintf("%s_%s_%s", prefix, path.Base(mediaPath), utils.GenerateToken())
	mediaName = models.SanitizeMediaName(mediaName)
	mediaName = mediaName + extension
	return mediaName
}

func generateUniqueMediaName(mediaPath string) string {

	filename := path.Base(mediaPath)
	baseName := filename[0 : len(filename)-len(path.Ext(filename))]
	baseExt := path.Ext(filename)

	mediaName := fmt.Sprintf("%s_%s", baseName, utils.GenerateToken())
	mediaName = models.SanitizeMediaName(mediaName) + baseExt

	return mediaName
}

func saveOriginalPhotoToDB(photo *models.Media, imageData *media_encoding.EncodeMediaData, photoDimensions *media_utils.PhotoDimensions) (*models.MediaURL, error) {
	originalImageName := generateUniqueMediaName(photo.Path)

	contentType, err := imageData.ContentType()
	if err != nil {
		return nil, err
	}

	fileStats, err := os.Stat(photo.Path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file stats of original photo")
	}

	mediaURL := models.MediaURL{
		Media:       photo,
		MediaName:   originalImageName,
		Width:       photoDimensions.Width,
		Height:      photoDimensions.Height,
		Purpose:     models.MediaOriginal,
		ContentType: string(*contentType),
		FileSize:    fileStats.Size(),
	}
	sql_media_urls_in := fmt.Sprintf("INSERT INTO `media_urls` (`created_at`,`updated_at`,`media_id`,`media_name`,`width`,`height`,`purpose`,`content_type`,`file_size`) VALUES (NOW(),NOW(),%v,'%v',%v,%v,'%v','%v',%v)", mediaURL.Media.ID, mediaURL.MediaName, mediaURL.Width, mediaURL.Height, mediaURL.Purpose, mediaURL.ContentType, mediaURL.FileSize)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_media_urls_in)
	return &mediaURL, nil
}
