package actions

//修改完
import (
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
	"strings"
	"time"
)

//修改完，未测试
func MyAlbums(user *models.User, order *models.Ordering, paginate *models.Pagination, onlyRoot *bool, showEmpty *bool, onlyWithFavorites *bool) ([]*models.Album, error) {
	if err := user.FillAlbums(); err != nil {
		return nil, err
	}

	if len(user.Albums) == 0 {
		return nil, nil
	}

	userAlbumIDs := make([]int, len(user.Albums))
	for i, album := range user.Albums {
		userAlbumIDs[i] = album.ID
	}
	var sql_albums_se string
	//query := db.Model(models.Album{}).Where("id IN (?)", userAlbumIDs)
	userAlbumID, _ := json.Marshal(userAlbumIDs)
	userAlbumids := strings.Trim(string(userAlbumID), "[]")
	sql_albums_se = fmt.Sprintf("select * from albums where id in (%v)", userAlbumids)

	if onlyRoot != nil && *onlyRoot {

		var singleRootAlbumID int = -1
		for _, album := range user.Albums {
			if album.ParentAlbumID == nil {
				if singleRootAlbumID == -1 {
					singleRootAlbumID = album.ID
				} else {
					singleRootAlbumID = -1
					break
				}
			}
		}

		if singleRootAlbumID != -1 && len(user.Albums) > 1 {
			//query = query.Where("parent_album_id = ?", singleRootAlbumID)
			sql_albums_se = sql_albums_se + fmt.Sprintf(" and parent_album_id = %v", singleRootAlbumID)
		} else {
			//query = query.Where("parent_album_id IS NULL")
			sql_albums_se += fmt.Sprintf(" and parent_album_id IS NULL")
		}
	}

	//这里注意一下
	if showEmpty == nil || !*showEmpty {
		//subQuery := db.Model(&models.Media{}).Where("album_id = albums.id")
		sub := fmt.Sprintf("select * from media where album_id = albums.id")
		if onlyWithFavorites != nil && *onlyWithFavorites {
			//favoritesSubquery := db.
			//	Model(&models.UserMediaData{UserID: user.ID}).
			//	Where("user_media_data.media_id = media.id").
			//	Where("user_media_data.favorite = true")
			favoritesSub := fmt.Sprintf("select * from user_media_data where user_id= %v and user_media_data.media_id = media.id and user_media_data.favorite = true", user.ID)
			//subQuery = subQuery.Where("EXISTS (?)", favoritesSubquery)
			sub += fmt.Sprintf(" and EXISTS (%v)", favoritesSub)
		}

		//query = query.Where("EXISTS (?)", subQuery)
		sql_albums_se += fmt.Sprintf(" and EXISTS (%v)", sub)
	}

	//query = models.FormatSQL(query, order, paginate)
	var orderby string
	orderby = *order.OrderBy
	sql_albums_se += fmt.Sprintf(" order by %v", orderby)

	var albums []*models.Album
	//if err := query.Find(&albums).Error; err != nil { //SELECT * FROM `albums` WHERE id IN (1) AND parent_album_id IS NULL ORDER BY `title`
	//	return nil, err
	//}
	dataApi, _ := DataApi.NewDataApiClient()
	res, _ := dataApi.Query(sql_albums_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var album models.Album
		album.ID = DataApi.GetInt(res, i, 0)
		album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, i, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
		album.Path = DataApi.GetString(res, i, 5)
		album.PathHash = DataApi.GetString(res, i, 6)
		album.CoverID = DataApi.GetIntP(res, i, 7)
		albums = append(albums, &album)
	}
	return albums, nil
}

//修改完
func Album(user *models.User, id int) (*models.Album, error) {
	var album models.Album
	//if err := db.First(&album, id).Error; err != nil { //SELECT * FROM `albums` WHERE `albums`.`id` = 1 ORDER BY `albums`.`id` LIMIT 1
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, errors.New("album not found")
	//	}
	//	return nil, err
	//}
	sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE `albums`.`id` = %v ORDER BY `albums`.`id` LIMIT 1", id)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	if len(res) == 0 {
		return nil, err
	}
	album.ID = DataApi.GetInt(res, 0, 0)
	album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	album.Title = DataApi.GetString(res, 0, 3)
	album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
	album.Path = DataApi.GetString(res, 0, 5)
	album.PathHash = DataApi.GetString(res, 0, 6)
	album.CoverID = DataApi.GetIntP(res, 0, 7)
	ownsAlbum, err := user.OwnsAlbum(&album) //这里注意一下
	if err != nil {
		return nil, err
	}

	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	return &album, nil
}

//修改完
func AlbumPath(user *models.User, album *models.Album) ([]*models.Album, error) {
	var album_path []*models.Album

	//err := db.Raw(`
	//	WITH recursive path_albums AS (
	//		SELECT * FROM albums anchor WHERE anchor.id = ?
	//		UNION
	//		SELECT parent.* FROM path_albums child JOIN albums parent ON parent.id = child.parent_album_id
	//	)
	//	SELECT * FROM path_albums WHERE id != ?
	//`, album.ID, album.ID).Scan(&album_path).Error
	/* WITH recursive path_albums AS (
	      SELECT * FROM albums anchor WHERE anchor.id = 1
	      UNION
	      SELECT parent.* FROM path_albums child JOIN albums parent ON parent.id = child.parent_album_id
	)
	SELECT * FROM path_albums WHERE id != 1
	*/

	sql_albums_se := fmt.Sprintf("WITH recursive path_albums AS (SELECT * FROM albums anchor WHERE anchor.id = %v UNION SELECT parent.* FROM path_albums child JOIN albums parent ON parent.id = child.parent_album_id)SELECT * FROM path_albums WHERE id != %v", album.ID, album.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_albums_se)
	num := len(res)
	for i := 0; i < num; i++ {
		var Album models.Album
		Album.ID = DataApi.GetInt(res, i, 0)
		Album.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
		Album.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
		Album.Title = DataApi.GetString(res, i, 3)
		Album.ParentAlbumID = DataApi.GetIntP(res, i, 4)
		Album.Path = DataApi.GetString(res, i, 5)
		Album.PathHash = DataApi.GetString(res, i, 6)
		Album.CoverID = DataApi.GetIntP(res, i, 7)
		album_path = append(album_path, &Album)
	}
	// Make sure to only return albums this user owns
	for i := len(album_path) - 1; i >= 0; i-- {
		album := album_path[i]

		owns, err := user.OwnsAlbum(album) //这里注意一下
		if err != nil {
			return nil, err
		}

		if !owns {
			album_path = album_path[i+1:]
			break
		}

	}

	if err != nil {
		return nil, err
	}

	return album_path, nil
}

//修改完
func SetAlbumCover(user *models.User, mediaID int) (*models.Album, error) {
	var media models.Media
	dataApi, _ := DataApi.NewDataApiClient()
	//if err := db.Find(&media, mediaID).Error; err != nil {
	//return nil, err
	//}
	sql_media_se := fmt.Sprintf("select * from media where media.id=%v", mediaID)
	res, err := dataApi.Query(sql_media_se)
	if len(res) == 0 {
		return nil, err
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

	var album models.Album

	//if err := db.Find(&album, &media.AlbumID).Error; err != nil {
	//	return nil, err
	//}
	sql_albums_se := fmt.Sprintf("select * from albums where id=%v", media.AlbumID)
	res, err = dataApi.Query(sql_albums_se)
	if len(res) == 0 {
		return nil, err
	}
	album.ID = int(*res[0][0].LongValue)
	album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	album.Title = *res[0][3].StringValue
	album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
	album.Path = *res[0][5].StringValue
	album.PathHash = *res[0][6].StringValue
	album.CoverID = DataApi.GetIntP(res, 0, 7)

	ownsAlbum, err := user.OwnsAlbum(&album) //这里注意一下，还未修改
	if err != nil {
		return nil, err
	}

	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	//if err := db.Model(&album).Update("cover_id", mediaID).Error; err != nil {
	//	return nil, err
	//}
	sql_albums_up := fmt.Sprintf("update albums set cover_id=%v where id =%v", mediaID, album.ID)
	dataApi.ExecuteSQl(sql_albums_up)
	album.CoverID = &mediaID
	return &album, nil
}

//修改完
func ResetAlbumCover(user *models.User, albumID int) (*models.Album, error) {
	var album models.Album
	//if err := db.Find(&album, albumID).Error; err != nil {
	//	return nil, err
	//}
	dataApi, _ := DataApi.NewDataApiClient()
	sql_albums_se := fmt.Sprintf("select * from albums where id=%v", albumID)
	res, err := dataApi.Query(sql_albums_se)
	album.ID = int(*res[0][0].LongValue)
	album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
	album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
	album.Title = *res[0][3].StringValue
	album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
	album.Path = *res[0][5].StringValue
	album.PathHash = *res[0][6].StringValue
	album.CoverID = DataApi.GetIntP(res, 0, 7)
	ownsAlbum, err := user.OwnsAlbum(&album) //这里注意一下还未修改
	if err != nil {
		return nil, err
	}

	if !ownsAlbum {
		return nil, errors.New("forbidden")
	}

	//if err := db.Model(&album).Update("cover_id", nil).Error; err != nil {
	//	return nil, err
	//}
	sql_albums_up := fmt.Sprintf("update albums set cover_id=NULL where id =%v", albumID)
	dataApi.ExecuteSQl(sql_albums_up)
	album.CoverID = nil
	return &album, nil
}
