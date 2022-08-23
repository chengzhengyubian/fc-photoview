package processing_tasks

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"os"
	"path"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/media_encoding/media_utils"
	"github.com/pkg/errors"
)

func generateSaveHighResJPEG(media *models.Media, imageData *media_encoding.EncodeMediaData, highres_name string, imagePath string, mediaURL *models.MediaURL) (*models.MediaURL, error) {

	err := imageData.EncodeHighRes(imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "creating high-res cached image")
	}

	photoDimensions, err := media_utils.GetPhotoDimensions(imagePath)
	if err != nil {
		return nil, err
	}

	fileStats, err := os.Stat(imagePath)
	if err != nil {
		return nil, errors.Wrap(err, "reading file stats of highres photo")
	}

	if mediaURL == nil {

		mediaURL = &models.MediaURL{
			MediaID:     media.ID,
			MediaName:   highres_name,
			Width:       photoDimensions.Width,
			Height:      photoDimensions.Height,
			Purpose:     models.PhotoHighRes,
			ContentType: "image/jpeg",
			FileSize:    fileStats.Size(),
		}

		sql_media_urls_in := fmt.Sprintf("INSERT INTO `media_urls` (`created_at`,`updated_at`,`media_id`,`media_name`,`width`,`height`,`purpose`,`content_type`,`file_size`) VALUES (NOW(),NOW(),%v,'%v',%v,%v,'%v','%v',%v)", mediaURL.MediaID, mediaURL.MediaName, mediaURL.Width, mediaURL.Height, mediaURL.Purpose, mediaURL.ContentType, mediaURL.FileSize)
		dataApi, _ := DataApi.NewDataApiClient()
		dataApi.ExecuteSQl(sql_media_urls_in)
	} else {
		mediaURL.Width = photoDimensions.Width
		mediaURL.Height = photoDimensions.Height
		mediaURL.FileSize = fileStats.Size()
		sql_media_urls_up := fmt.Sprintf("update media_urls set updated_at=NOW(),width=%v,`height`=%v,file_size=%v", mediaURL.Width, mediaURL.Height, mediaURL.FileSize)
		dataApi, _ := DataApi.NewDataApiClient()
		dataApi.ExecuteSQl(sql_media_urls_up)
	}

	return mediaURL, nil
}

//修改完，测试基本成功，入库有延迟
func generateSaveThumbnailJPEG(media *models.Media, thumbnail_name string, photoCachePath string, baseImagePath string, mediaURL *models.MediaURL) (*models.MediaURL, error) {
	thumbOutputPath := path.Join(photoCachePath, thumbnail_name)

	thumbSize, err := media_encoding.EncodeThumbnail(baseImagePath, thumbOutputPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not create thumbnail cached image")
	}

	fileStats, err := os.Stat(thumbOutputPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading file stats of thumbnail photo")
	}

	if mediaURL == nil {

		mediaURL = &models.MediaURL{
			MediaID:     media.ID,
			MediaName:   thumbnail_name,
			Width:       thumbSize.Width,
			Height:      thumbSize.Height,
			Purpose:     models.PhotoThumbnail,
			ContentType: "image/jpeg",
			FileSize:    fileStats.Size(),
		}
		sql_media_urls_in := fmt.Sprintf("INSERT INTO `media_urls` (`created_at`,`updated_at`,`media_id`,`media_name`,`width`,`height`,`purpose`,`content_type`,`file_size`) VALUES (NOW(),NOW(),%v,'%v',%v,%v,'%v','%v',%v)", mediaURL.MediaID, mediaURL.MediaName, mediaURL.Width, mediaURL.Height, mediaURL.Purpose, mediaURL.ContentType, mediaURL.FileSize)
		dataApi, _ := DataApi.NewDataApiClient()
		dataApi.ExecuteSQl(sql_media_urls_in)
	} else {
		mediaURL.Width = thumbSize.Width
		mediaURL.Height = thumbSize.Height
		mediaURL.FileSize = fileStats.Size()
		sql_media_urls_up := fmt.Sprintf("update media_urls set updated_at=NOW(),width=%v,`height`=%v,file_size=%v", mediaURL.Width, mediaURL.Height, mediaURL.FileSize)
		dataApi, _ := DataApi.NewDataApiClient()
		dataApi.ExecuteSQl(sql_media_urls_up)
	}
	return mediaURL, nil
}
