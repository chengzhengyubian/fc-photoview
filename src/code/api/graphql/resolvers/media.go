package resolvers

//修改完
import (
	"context"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strconv"
	"strings"
	"time"

	"github.com/photoview/photoview/api/dataloader"
	api "github.com/photoview/photoview/api/graphql"
	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/pkg/errors"
)

func (r *queryResolver) MyMedia(ctx context.Context, order *models.Ordering, paginate *models.Pagination) ([]*models.Media, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	return actions.MyMedia(user, order, paginate)
}

//修改完，未测试
func (r *queryResolver) Media(ctx context.Context, id int, tokenCredentials *models.ShareTokenCredentials) (*models.Media, error) {
	//db := r.DB(ctx)
	if tokenCredentials != nil {

		shareToken, err := r.ShareToken(ctx, *tokenCredentials)
		if err != nil {
			return nil, err
		}

		if *shareToken.MediaID == id {
			return shareToken.Media, nil
		}
	}
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}
	var media models.Media
	//err := db.
	//	Joins("Album").
	//	Where("media.id = ?", id).
	//	Where("EXISTS (SELECT * FROM user_albums WHERE user_albums.album_id = media.album_id AND user_albums.user_id = ?)", user.ID).
	//	Where("media.id IN (?)", db.Model(&models.MediaURL{}).Select("media_id").Where("media_urls.media_id = media.id")).
	//	First(&media).Error
	subsql := fmt.Sprintf("select media_id from media_urls where media_urls.media_id = media.id")
	sql_album_se := fmt.Sprintf("select media.*,albums.* from media left join albums on media.album_id=albums.id where media.id= %v and EXISTS (SELECT * FROM user_albums WHERE user_albums.album_id = media.album_id AND user_albums.user_id = %v) and media.id IN(%v) limit 1", id, user.ID, subsql)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_album_se)
	if len(res) == 0 {
		return nil, errors.Wrap(err, "could not get media by media_id and user_id from database")
	}
	media.ID = DataApi.GetInt(res, 0, 0)
	media.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	media.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	media.Title = *res[0][3].StringValue
	media.Path = *res[0][4].StringValue
	media.PathHash = *res[0][5].StringValue
	media.AlbumID = int(*res[0][6].LongValue)
	media.ExifID = DataApi.GetIntP(res, 0, 7)
	media.DateShot = time.Unix(*res[0][8].LongValue/1000, 0)
	if *res[0][9].StringValue == "photo" {
		media.Type = models.MediaTypePhoto
	} else {
		media.Type = models.MediaTypeVideo
	}
	media.VideoMetadataID = DataApi.GetIntP(res, 0, 10)
	media.SideCarPath = DataApi.GetStringP(res, 0, 11)
	media.SideCarHash = DataApi.GetStringP(res, 0, 12)
	media.Blurhash = DataApi.GetStringP(res, 0, 13)
	media.AlbumID = DataApi.GetInt(res, 0, 14)
	media.Album.ID = DataApi.GetInt(res, 0, 14)
	media.Album.CreatedAt = time.Unix(*res[0][15].LongValue/1000, 0)
	media.Album.UpdatedAt = time.Unix(*res[0][16].LongValue/1000, 0)
	media.Album.Title = DataApi.GetString(res, 0, 17)
	media.Album.ParentAlbumID = DataApi.GetIntP(res, 0, 18)
	media.Album.Path = DataApi.GetString(res, 0, 19)
	media.Album.PathHash = DataApi.GetString(res, 0, 20)
	media.Album.CoverID = DataApi.GetIntP(res, 0, 21)
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not get media by media_id and user_id from database")
	//}
	return &media, nil
}

//修改完，未测试
func (r *queryResolver) MediaList(ctx context.Context, ids []int) ([]*models.Media, error) {
	//db := r.DB(ctx)
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	if len(ids) == 0 {
		return nil, errors.New("no ids provided")
	}

	var media []*models.Media
	//err := db.Model(&media).
	//	Joins("LEFT JOIN user_albums ON user_albums.album_id = media.album_id").
	//	Where("media.id IN ?", ids).
	//	Where("user_albums.user_id = ?", user.ID).
	//	Find(&media).Error
	//
	//if err != nil {
	//	return nil, errors.Wrap(err, "could not get media list by media_id and user_id from database")
	//}

	id, _ := json.Marshal(ids)
	IDs := strings.Trim(string(id), "[]")
	sql_media_se := fmt.Sprintf("select media.* from media left join user_albums ON user_albums.album_id = media.album_id where media_id in (%v) and user_albums.user_id =%v", IDs, user.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, _ := dataApi.Query(sql_media_se)
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
	return media, nil
}

type mediaResolver struct {
	*Resolver
}

func (r *Resolver) Media() api.MediaResolver {
	return &mediaResolver{r}
}

func (r *mediaResolver) Type(ctx context.Context, media *models.Media) (models.MediaType, error) {
	formattedType := models.MediaType(strings.Title(string(media.Type)))
	return formattedType, nil
}

//修改完
func (r *mediaResolver) Album(ctx context.Context, obj *models.Media) (*models.Album, error) {
	var album models.Album
	//err := r.DB(ctx).Find(&album, obj.AlbumID).Error // SELECT * FROM `albums` WHERE `albums`.`id` = 1
	sql_albums_se := "SELECT * FROM `albums` WHERE `albums`.`id` =" + strconv.Itoa(obj.AlbumID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	album.ID = DataApi.GetInt(res, 0, 0)
	album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	album.Title = DataApi.GetString(res, 0, 3)
	album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
	album.Path = DataApi.GetString(res, 0, 5)
	album.PathHash = DataApi.GetString(res, 0, 6)
	album.CoverID = DataApi.GetIntP(res, 0, 7)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

//修改完，未测试
func (r *mediaResolver) Shares(ctx context.Context, media *models.Media) ([]*models.ShareToken, error) {
	var shareTokens []*models.ShareToken
	//if err := r.DB(ctx).Where("media_id = ?", media.ID).Find(&shareTokens).Error; err != nil {
	//	return nil, errors.Wrapf(err, "get shares for media (%s)", media.Path)
	//}
	sql_share_tokens_select := fmt.Sprintf("select * from share_tokens where media_id=%v", media.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_share_tokens_select)
	if len(res) == 0 {
		return nil, errors.Wrapf(err, "get shares for media (%s)", media.Path)
	}
	for i := 0; i < len(res); i++ {
		var ShareToken models.ShareToken
		ShareToken.ID = DataApi.GetInt(res, i, 0)
		ShareToken.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		ShareToken.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		ShareToken.Value = DataApi.GetString(res, i, 3)
		ShareToken.OwnerID = DataApi.GetInt(res, i, 4)
		if DataApi.GetLongP(res, 0, 5) == nil {
			ShareToken.Expire = nil
		} else {
			expire := time.Unix(DataApi.GetLong(res, 0, 5)/1000, 0)
			ShareToken.Expire = &expire
		}
		//expire := time.Unix(DataApi.GetLongP(res,0,5)/1000, 0)
		//ShareToken.Expire = &expire
		ShareToken.Password = DataApi.GetStringP(res, i, 6)
		ShareToken.AlbumID = DataApi.GetIntP(res, i, 7)
		ShareToken.MediaID = DataApi.GetIntP(res, i, 8)
		shareTokens = append(shareTokens, &ShareToken)
	}
	return shareTokens, nil
}

//修改完，未测试
func (r *mediaResolver) Downloads(ctx context.Context, media *models.Media) ([]*models.MediaDownload, error) {

	var mediaUrls []*models.MediaURL
	//if err := r.DB(ctx).Where("media_id = ?", media.ID).Find(&mediaUrls).Error; err != nil {
	//	return nil, errors.Wrapf(err, "get downloads for media (%s)", media.Path)
	//}
	sql_media_urls_select := fmt.Sprintf("select * from media_urls where media_id = %v", media.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_urls_select)
	if len(res) == 0 {
		return nil, errors.Wrapf(err, "get downloads for media (%s)", media.Path)
	}
	for i := 0; i < len(res); i++ {
		var mediaUrl models.MediaURL
		mediaUrl.ID = DataApi.GetInt(res, i, 0)
		mediaUrl.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		mediaUrl.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		mediaUrl.MediaID = DataApi.GetInt(res, i, 3)
		mediaUrl.MediaName = DataApi.GetString(res, i, 4)
		mediaUrl.Width = DataApi.GetInt(res, i, 5)
		mediaUrl.Height = DataApi.GetInt(res, i, 6)
		switch DataApi.GetString(res, i, 7) {
		case "thumbnail":
			mediaUrl.Purpose = models.PhotoThumbnail
		case "high-res":
			mediaUrl.Purpose = models.PhotoHighRes
		case "original":
			mediaUrl.Purpose = models.MediaOriginal
		case "video-web":
			mediaUrl.Purpose = models.VideoWeb
		case "video-thumbnail":
			mediaUrl.Purpose = models.VideoThumbnail
		}
		mediaUrl.ContentType = DataApi.GetString(res, i, 8)
		mediaUrl.FileSize = DataApi.GetLong(res, i, 9)
		mediaUrls = append(mediaUrls, &mediaUrl)
	}
	downloads := make([]*models.MediaDownload, 0)

	for _, url := range mediaUrls {

		var title string
		switch {
		case url.Purpose == models.MediaOriginal:
			title = "Original"
		case url.Purpose == models.PhotoThumbnail:
			title = "Small"
		case url.Purpose == models.PhotoHighRes:
			title = "Large"
		case url.Purpose == models.VideoThumbnail:
			title = "Video thumbnail"
		case url.Purpose == models.VideoWeb:
			title = "Web optimized video"
		}

		downloads = append(downloads, &models.MediaDownload{
			Title:    title,
			MediaURL: url,
		})
	}

	return downloads, nil
}

func (r *mediaResolver) HighRes(ctx context.Context, media *models.Media) (*models.MediaURL, error) {
	if media.Type != models.MediaTypePhoto {
		return nil, nil
	}

	return dataloader.For(ctx).MediaHighres.Load(media.ID)
}

func (r *mediaResolver) Thumbnail(ctx context.Context, media *models.Media) (*models.MediaURL, error) {
	return dataloader.For(ctx).MediaThumbnail.Load(media.ID)
}

func (r *mediaResolver) VideoWeb(ctx context.Context, media *models.Media) (*models.MediaURL, error) {
	if media.Type != models.MediaTypeVideo {
		return nil, nil
	}

	return dataloader.For(ctx).MediaVideoWeb.Load(media.ID)
}

//修改完，未测试
func (r *mediaResolver) Exif(ctx context.Context, media *models.Media) (*models.MediaEXIF, error) {
	if media.Exif != nil {
		return media.Exif, nil
	}
	var exif models.MediaEXIF
	//if err := r.DB(ctx).Model(&media).Association("Exif").Find(&exif); err != nil {
	//	return nil, err
	//}
	sql_media_exif_select := fmt.Sprintf("select media_exif.* from media_exif left join media on media.exif_id=media.exif_id and media.id=%v order by media.exif_id limit 1", media.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_exif_select)
	if len(res) == 0 {
		return nil, err
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

func (r *mediaResolver) Favorite(ctx context.Context, media *models.Media) (bool, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return false, auth.ErrUnauthorized
	}

	return dataloader.For(ctx).UserMediaFavorite.Load(&models.UserMediaData{
		UserID:  user.ID,
		MediaID: media.ID,
	})
}

func (r *mutationResolver) FavoriteMedia(ctx context.Context, mediaID int, favorite bool) (*models.Media, error) {

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return user.FavoriteMedia( /*r.DB(ctx), */ mediaID, favorite)
}

//修改完，未测试
func (r *mediaResolver) Faces(ctx context.Context, media *models.Media) ([]*models.ImageFace, error) {
	if face_detection.GlobalFaceDetector == nil {
		return []*models.ImageFace{}, nil
	}

	if media.Faces != nil {
		return media.Faces, nil
	}

	var faces []*models.ImageFace
	//if err := r.DB(ctx).Model(&media).Association("Faces").Find(&faces); err != nil {
	//	return nil, err
	//}

	sql_image_faces_select := fmt.Sprintf("select image_faces.* from image_faces left join media on image_faces.media_id=media.id where media.id=%v", media.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, _ := dataApi.Query(sql_image_faces_select)
	for i := 0; i < len(res); i++ {
		var face models.ImageFace
		face.ID = DataApi.GetInt(res, i, 0)
		face.CreatedAt = time.Unix(DataApi.GetLong(res, i, 1)/1000, 0)
		face.UpdatedAt = time.Unix(DataApi.GetLong(res, i, 2)/1000, 0)
		face.FaceGroupID = DataApi.GetInt(res, i, 3)
		face.MediaID = DataApi.GetInt(res, i, 4)
		//face.Descriptor=*res[i][5].ArrayValue
		//face.Rectangle=DataApi.GetStringP()
		faces = append(faces, &face)
	}
	return faces, nil
}
