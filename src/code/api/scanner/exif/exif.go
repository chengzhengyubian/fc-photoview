package exif

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"log"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
)

type ExifParser interface {
	ParseExif(media_path string) (*models.MediaEXIF, error)
}

var globalExifParser ExifParser

func InitializeEXIFParser() {
	exiftoolParser, err := NewExiftoolParser()

	if err != nil {
		log.Printf("Failed to get exiftool, using internal exif parser instead: %v\n", err)
		globalExifParser = NewInternalExifParser()
	} else {
		log.Println("Found exiftool")
		globalExifParser = exiftoolParser
	}
}

// SaveEXIF scans the media file for exif metadata and saves it in the database if found
func SaveEXIF( /*tx *gorm.DB, */ media *models.Media) (*models.MediaEXIF, error) {

	{
		if media.ExifID != nil {

			var exif models.MediaEXIF
			sql_media_exif_select := fmt.Sprintf("select * from media_exif where id=%v", media.ExifID)
			dataApi, _ := DataApi.NewDataApiClient()
			res, err := dataApi.Query(sql_media_exif_select)
			if len(res) == 0 {
				return nil, errors.Wrap(err, "get EXIF for media from database")
			}
			exif.ID = DataApi.GetInt(res, 0, 0)
			exif.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
			exif.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 2)/1000, 0)
			exif.Camera = DataApi.GetStringP(res, 0, 3)
			exif.Maker = DataApi.GetStringP(res, 0, 4)
			exif.Lens = DataApi.GetStringP(res, 0, 5)
			date := time.Unix(DataApi.GetLong(res, 0, 6)/1000, 0)
			exif.DateShot = &date
			exposure := float64(DataApi.GetLong(res, 0, 7))
			exif.Exposure = &exposure
			aperture := float64(DataApi.GetLong(res, 0, 8))
			exif.Aperture = &aperture
			exif.Iso = DataApi.GetLongP(res, 0, 9)
			focalLength := float64(DataApi.GetLong(res, 0, 10))
			exif.FocalLength = &focalLength
			exif.Flash = DataApi.GetLongP(res, 0, 11)
			exif.Orientation = DataApi.GetLongP(res, 0, 12)
			exif.ExposureProgram = DataApi.GetLongP(res, 0, 13)
			gpsLaitude := float64(DataApi.GetLong(res, 0, 14))
			exif.GPSLatitude = &gpsLaitude
			gpsLongtude := float64(DataApi.GetLong(res, 0, 15))
			exif.GPSLongitude = &gpsLongtude
			exif.Description = DataApi.GetStringP(res, 0, 16)
			return &exif, nil
		}
	}

	if globalExifParser == nil {
		return nil, errors.New("No exif parser initialized")
	}

	exif, err := globalExifParser.ParseExif(media.Path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse exif data")
	}

	if exif == nil {
		return nil, nil
	}
	if exif.DateShot != nil && !exif.DateShot.Equal(media.DateShot) {
		media.DateShot = *exif.DateShot
		timestr := media.DateShot.Format("2006-01-02 15:04:05")
		sql_media_up := fmt.Sprintf("update media set data_shot=%v where media.id=%v", timestr, media.ID)
		dataApi, _ := DataApi.NewDataApiClient()
		dataApi.ExecuteSQl(sql_media_up)
	}

	return exif, nil
}
