package routes

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
	"net/http"
	"time"
)

func authenticateMedia(media *models.Media /*db *gorm.DB,*/, r *http.Request) (success bool, responseMessage string, responseStatus int, errorMessage error) {
	user := auth.UserFromContext(r.Context())

	if user != nil {
		var album models.Album
		sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE `albums`.`id` =%v ORDER BY `albums`.`id` LIMIT 1", media.AlbumID)
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.Query(sql_albums_se)
		if len(res) == 0 {
			return false, "internal server error", http.StatusInternalServerError, err
		}
		album.ID = DataApi.GetInt(res, 0, 0)
		album.CreatedAt = time.Unix(*res[0][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res[0][2].LongValue/1000, 0)
		album.Title = DataApi.GetString(res, 0, 3)
		album.ParentAlbumID = DataApi.GetIntP(res, 0, 4)
		album.Path = DataApi.GetString(res, 0, 5)
		album.PathHash = DataApi.GetString(res, 0, 6)
		album.CoverID = DataApi.GetIntP(res, 0, 7)

		ownsAlbum, err := user.OwnsAlbum(&album) //这里注意一下，这里还没改
		if err != nil {
			return false, "internal server error", http.StatusInternalServerError, err
		}

		if !ownsAlbum {
			return false, "invalid credentials", http.StatusForbidden, nil
		}
	} else {
		if success, respMsg, respStatus, err := shareTokenFromRequest( /*db, */ r, &media.ID, &media.AlbumID); !success {
			return success, respMsg, respStatus, err
		}

	}

	return true, "success", http.StatusAccepted, nil
}

func authenticateAlbum(album *models.Album /*db *gorm.DB, */, r *http.Request) (success bool, responseMessage string, responseStatus int, errorMessage error) {
	user := auth.UserFromContext(r.Context())

	if user != nil {
		ownsAlbum, err := user.OwnsAlbum(album)
		if err != nil {
			return false, "internal server error", http.StatusInternalServerError, err
		}

		if !ownsAlbum {
			return false, "invalid credentials", http.StatusForbidden, nil
		}
	} else {
		if success, respMsg, respStatus, err := shareTokenFromRequest( /*db, */ r, nil, &album.ID); !success {
			return success, respMsg, respStatus, err
		}

	}

	return true, "success", http.StatusAccepted, nil
}

//修改完
func shareTokenFromRequest( /*db *gorm.DB,*/ r *http.Request, mediaID *int, albumID *int) (success bool, responseMessage string, responseStatus int, errorMessage error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return false, "unauthorized", http.StatusForbidden, errors.New("share token not provided")
	}

	var shareToken models.ShareToken
	sql_share_tokens_se := "select * from  share_tokens where share_tokens.value=\"" + token + "\" order by share_tokens.id limit 1"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_share_tokens_se)
	if err != nil {
		return false, "internal server error", http.StatusInternalServerError, err
	}
	shareToken.ID = DataApi.GetInt(res, 0, 0)
	shareToken.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
	shareToken.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 2)/1000, 0)
	shareToken.Value = DataApi.GetString(res, 0, 3)
	shareToken.OwnerID = DataApi.GetInt(res, 0, 4)
	Time := time.Unix(DataApi.GetLong(res, 0, 5)/1000, 0)
	shareToken.Expire = &Time
	shareToken.Password = DataApi.GetStringP(res, 0, 6)
	shareToken.AlbumID = DataApi.GetIntP(res, 0, 7)
	shareToken.MediaID = DataApi.GetIntP(res, 0, 8)
	// Validate share token password, if set
	if shareToken.Password != nil {
		tokenPasswordCookie, err := r.Cookie(fmt.Sprintf("share-token-pw-%s", shareToken.Value))
		if err != nil {
			return false, "unauthorized", http.StatusForbidden, errors.Wrap(err, "get share token password cookie")
		}
		// tokenPassword := r.Header.Get("TokenPassword")
		tokenPassword := tokenPasswordCookie.Value

		if err := bcrypt.CompareHashAndPassword([]byte(*shareToken.Password), []byte(tokenPassword)); err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return false, "unauthorized", http.StatusForbidden, errors.New("incorrect password for share token")
			} else {
				return false, "internal server error", http.StatusInternalServerError, err
			}
		}
	}

	if shareToken.AlbumID != nil && albumID == nil {
		return false, "unauthorized", http.StatusForbidden, errors.New("share token is of type album, but no albumID was provided to function")
	}

	if shareToken.MediaID != nil && mediaID == nil {
		return false, "unauthorized", http.StatusForbidden, errors.New("share token is of type media, but no mediaID was provided to function")
	}

	if shareToken.AlbumID != nil && *albumID != *shareToken.AlbumID {

		var count int
		sql_albums_se := fmt.Sprintf("WITH recursive child_albums AS (SELECT * FROM albums WHERE parent_album_id = %v UNION ALL SELECT child.* FROM albums child JOIN child_albums parent ON parent.id = child.parent_album_id)SELECT COUNT(id) FROM child_albums WHERE id =%v", *shareToken.AlbumID, albumID)
		res, err = dataApi.Query(sql_albums_se)
		if err != nil {
			return false, "internal server error", http.StatusInternalServerError, err
		}
		count = DataApi.GetInt(res, 0, 0)
		if count == 0 {
			return false, "unauthorized", http.StatusForbidden, errors.New("no child albums found for share token")
		}
	}

	if shareToken.MediaID != nil && *mediaID != *shareToken.MediaID {
		return false, "unauthorized", http.StatusForbidden, errors.New("media share token does not match mediaID")
	}
	return true, "", 0, nil
}
