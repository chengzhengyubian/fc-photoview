package routes

//修改完
import (
	"archive/zip"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/photoview/photoview/api/graphql/models"
)

//修改完，待测试
func RegisterDownloadRoutes( /*db *gorm.DB,*/ router *mux.Router) {
	router.HandleFunc("/album/{album_id}/{media_purpose}", func(w http.ResponseWriter, r *http.Request) {
		albumID := mux.Vars(r)["album_id"]
		mediaPurpose := mux.Vars(r)["media_purpose"]
		mediaPurposeList := strings.SplitN(mediaPurpose, ",", 10)

		var album models.Album
		/*if err := db.Find(&album, albumID).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404"))
			return
		}*/
		sql_albums_se := fmt.Sprintf("select * from album where id=%v", albumID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, _ := dataApi.Query(sql_albums_se)
		album.ID = DataApi.GetInt(res, 0, 0)
		album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, 0, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
		album.Path = DataApi.GetString(res, 0, 5)
		album.PathHash = DataApi.GetString(res, 0, 6)
		album.CoverID = DataApi.GetIntP(res, 0, 7)
		if success, response, status, err := authenticateAlbum(&album /* db,*/, r); !success {
			if err != nil {
				log.Printf("WARN: error authenticating album for download: %v\n", err)
			}
			w.WriteHeader(status)
			w.Write([]byte(response))
			return
		}

		//var mediaWhereQuery string
		/*if drivers.POSTGRES.MatchDatabase(db) {
			mediaWhereQuery = "\"Media\".album_id = ?"
		} else {
			mediaWhereQuery = "Media.album_id = ?"
		}*/

		var mediaURLs []*models.MediaURL
		/*if err := db.Joins("Media").Where(mediaWhereQuery, album.ID).Where("media_urls.purpose IN (?)", mediaPurposeList).Find(&mediaURLs).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}*/
		mediapurpose, _ := json.Marshal(mediaPurposeList)
		mediapurposelist := strings.Trim(string(mediapurpose), "[]")
		/*SELECT `media_urls`.`id`,`media_urls`.`created_at`,`media_urls`.`updated_at`,`media_urls`.`media_id`,`media_urls`.`media_name`,`media_urls`.`width`,`media_urls`.`height`,`media_urls`.`purpose`,`media_urls`.`content_type`,`media_urls`.`file_size`,`Media`.`id` AS `Media__id`,`Media`.`created_at` AS `Media__created_at`,`Media`.`updated_at` AS `Media__updated_at`,`Media`.`title` AS `Media__title`,`Media`.`path` AS `Media__path`,`Media`.`path_hash` AS `Media__path_hash`,`Media`.`album_id` AS `Media__album_id`,`Media`.`exif_id` AS `Media__exif_id`,`Media`.`date_shot` AS `Media__date_shot`,`Media`.`type` AS `Media__type`,`Media`.`video_metadata_id` AS `Media__video_metadata_id`,`Media`.`side_car_path` AS `Media__side_car_path`,`Media`.`side_car_hash` AS `Media__side_car_hash`,`Media`.`blurhash` AS `Media__blurhash` FROM `media_urls` LEFT JOIN `media` `Media` ON `media_urls`.`media_id` = `Media`.`id` WHERE Media.album_id = 146 AND media_urls.purpose IN ('original')*/
		sql_media_urls_se := fmt.Sprintf("SELECT `media_urls`.`id`,`media_urls`.`created_at`,`media_urls`.`updated_at`,`media_urls`.`media_id`,`media_urls`.`media_name`,`media_urls`.`width`,`media_urls`.`height`,`media_urls`.`purpose`,`media_urls`.`content_type`,`media_urls`.`file_size`,`Media`.`id` AS `Media__id`,`Media`.`created_at` AS `Media__created_at`,`Media`.`updated_at` AS `Media__updated_at`,`Media`.`title` AS `Media__title`,`Media`.`path` AS `Media__path`,`Media`.`path_hash` AS `Media__path_hash`,`Media`.`album_id` AS `Media__album_id`,`Media`.`exif_id` AS `Media__exif_id`,`Media`.`date_shot` AS `Media__date_shot`,`Media`.`type` AS `Media__type`,`Media`.`video_metadata_id` AS `Media__video_metadata_id`,`Media`.`side_car_path` AS `Media__side_car_path`,`Media`.`side_car_hash` AS `Media__side_car_hash`,`Media`.`blurhash` AS `Media__blurhash` FROM `media_urls` LEFT JOIN `media` `Media` ON `media_urls`.`media_id` = `Media`.`id` WHERE Media.album_id = %v AND media_urls.purpose IN (%v)", album.ID, mediapurposelist)
		dataAPi, _ := DataApi.NewDataApiClient()
		res, _ = dataAPi.Query(sql_media_urls_se)
		if len(res) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
		//fmt.Print(res)
		num := len(res)
		for i := 0; i < num; i++ {
			var mediaURL models.MediaURL
			mediaURL.ID = int(*res[i][0].LongValue)
			mediaURL.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
			mediaURL.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
			mediaURL.MediaID = int(*res[i][3].LongValue)
			mediaURL.MediaName = *res[i][4].StringValue
			mediaURL.Width = int(*res[i][5].LongValue)
			mediaURL.Height = int(*res[i][6].LongValue)
			switch *res[i][7].StringValue {
			case "thumbnail":
				mediaURL.Purpose = models.PhotoThumbnail
			case "high-res":
				mediaURL.Purpose = models.PhotoHighRes
			case "original":
				mediaURL.Purpose = models.MediaOriginal
			case "video-web":
				mediaURL.Purpose = models.VideoWeb
			case "video-thumbnail":
				mediaURL.Purpose = models.VideoThumbnail
			}
			mediaURL.ContentType = *res[i][8].StringValue
			mediaURL.FileSize = *res[i][9].LongValue
			Media := models.Media{}
			mediaURL.Media = &Media
			//println("1\n")
			mediaURL.Media.ID = int(*res[i][10].LongValue)
			mediaURL.Media.CreatedAt = time.Unix(DataApi.GetLong(res, i, 11)/1000, 0)
			mediaURL.Media.UpdatedAt = time.Unix(DataApi.GetLong(res, i, 12)/1000, 0)
			mediaURL.Media.Title = DataApi.GetString(res, i, 13)
			mediaURL.Media.Path = DataApi.GetString(res, i, 14)
			mediaURL.Media.PathHash = DataApi.GetString(res, i, 15)
			mediaURL.Media.AlbumID = DataApi.GetInt(res, i, 16)
			mediaURL.Media.ExifID = DataApi.GetIntP(res, i, 17)
			mediaURL.Media.DateShot = time.Unix(DataApi.GetLong(res, i, 18)/1000, 0)
			switch DataApi.GetString(res, i, 19) {
			case "photo":
				mediaURL.Media.Type = models.MediaTypePhoto
			case "video":
				mediaURL.Media.Type = models.MediaTypeVideo
			}
			mediaURL.Media.VideoMetadataID = DataApi.GetIntP(res, i, 20)
			mediaURL.Media.SideCarPath = DataApi.GetStringP(res, i, 21)
			mediaURL.Media.SideCarHash = DataApi.GetStringP(res, i, 22)
			mediaURL.Media.Blurhash = DataApi.GetStringP(res, i, 23)
			mediaURLs = append(mediaURLs, &mediaURL)
		}

		if len(mediaURLs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no media found"))
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", album.Title))

		zipWriter := zip.NewWriter(w)

		for _, media := range mediaURLs {
			zipFile, err := zipWriter.Create(fmt.Sprintf("%s/%s", album.Title, media.MediaName))
			if err != nil {
				log.Printf("ERROR: Failed to create a file in zip, when downloading album (%d): %v\n", album.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}

			filePath, err := media.CachedPath()
			if err != nil {
				log.Printf("ERROR: Failed to get mediaURL cache path, when downloading album (%d): %v\n", album.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}

			fileData, err := os.Open(filePath)
			if err != nil {
				log.Printf("ERROR: Failed to open file to include in zip, when downloading album (%d): %v\n", album.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}

			_, err = io.Copy(zipFile, fileData)
			if err != nil {
				log.Printf("ERROR: Failed to copy file data, when downloading album (%d): %v\n", album.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}

			if err := fileData.Close(); err != nil {
				log.Printf("ERROR: Failed to close file, when downloading album (%d): %v\n", album.ID, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
				return
			}
		}

		// close the zip Writer to flush the contents to the ResponseWriter
		zipWriter.Close()
	})
}
