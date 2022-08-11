package resolvers

//修改完
import (
	"context"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"time"

	api "github.com/photoview/photoview/api/graphql"
	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

//
type shareTokenResolver struct {
	*Resolver
}

func (r *Resolver) ShareToken() api.ShareTokenResolver {
	return &shareTokenResolver{r}
}

func (r *shareTokenResolver) Owner(ctx context.Context, obj *models.ShareToken) (*models.User, error) {
	return &obj.Owner, nil
}

func (r *shareTokenResolver) Album(ctx context.Context, obj *models.ShareToken) (*models.Album, error) {
	return obj.Album, nil
}

func (r *shareTokenResolver) Media(ctx context.Context, obj *models.ShareToken) (*models.Media, error) {
	return obj.Media, nil
}

func (r *shareTokenResolver) HasPassword(ctx context.Context, obj *models.ShareToken) (bool, error) {
	hasPassword := obj.Password != nil
	return hasPassword, nil
}

//未修改
func (r *queryResolver) ShareToken(ctx context.Context, credentials models.ShareTokenCredentials) (*models.ShareToken, error) {

	var token models.ShareToken
	//if err := r.DB(ctx).Preload(clause.Associations).Where("value = ?", credentials.Token).First(&token).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, errors.New("share not found")
	//	} else {
	//		return nil, errors.Wrap(err, "failed to get share token from database")
	//	}
	//}

	if token.Password != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(*token.Password), []byte(*credentials.Password)); err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, errors.New("unauthorized")
			} else {
				return nil, errors.Wrap(err, "failed to compare token password hashes")
			}
		}
	}
	sql_share_tokens_se := fmt.Sprintf("select share_tokens.*,users.*,albums.*,media.* from share_tokens left join users on share_tokens.owner_id=users.id left join albums on share_tokens.album_id=albums.id left join media on share_tokens.media_id=media.id where share_tokens.value='%v'", credentials.Token)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_share_tokens_se)
	if len(res) == 0 {
		return nil, errors.Wrap(err, "share not found")
	}
	token.ID = DataApi.GetInt(res, 0, 0)
	token.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
	token.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 2)/1000, 0)
	token.Value = DataApi.GetString(res, 0, 3)
	token.OwnerID = DataApi.GetInt(res, 0, 4)
	expireS := time.Unix(DataApi.GetLong(res, 0, 5)/1000, 0)
	token.Expire = &expireS
	token.Password = DataApi.GetStringP(res, 0, 6)
	token.AlbumID = DataApi.GetIntP(res, 0, 7)
	token.MediaID = DataApi.GetIntP(res, 0, 8)
	token.Owner.ID = DataApi.GetInt(res, 0, 9)
	token.Owner.Username = DataApi.GetString(res, 0, 13)
	token.Owner.Password = DataApi.GetStringP(res, 0, 14)
	token.Owner.Admin = DataApi.GetBoolean(res, 0, 15)
	token.Album.ID = DataApi.GetInt(res, 0, 16)
	token.Album.CreatedAt = time.Unix(*res[0][17].LongValue/1000, 0)
	token.Album.UpdatedAt = time.Unix(*res[0][18].LongValue/1000, 0)
	token.Album.Title = DataApi.GetString(res, 0, 19)
	token.Album.ParentAlbumID = DataApi.GetIntP(res, 0, 20)
	token.Album.Path = DataApi.GetString(res, 0, 21)
	token.Album.PathHash = DataApi.GetString(res, 0, 22)
	token.Album.CoverID = DataApi.GetIntP(res, 0, 23)
	token.Media.ID = DataApi.GetInt(res, 0, 24)
	token.Media.CreatedAt = time.Unix(*res[0][25].LongValue/1000, 0)
	token.Media.UpdatedAt = time.Unix(*res[0][26].LongValue/1000, 0)
	token.Media.Title = *res[0][27].StringValue
	token.Media.Path = *res[0][28].StringValue
	token.Media.PathHash = *res[0][29].StringValue
	token.Media.AlbumID = int(*res[0][30].LongValue)
	token.Media.ExifID = DataApi.GetIntP(res, 0, 31)
	token.Media.DateShot = time.Unix(*res[0][32].LongValue/1000, 0)
	if *res[0][9].StringValue == "photo" {
		token.Media.Type = models.MediaTypePhoto
	} else {
		token.Media.Type = models.MediaTypeVideo
	}
	token.Media.VideoMetadataID = DataApi.GetIntP(res, 0, 33)
	token.Media.SideCarPath = DataApi.GetStringP(res, 0, 34)
	token.Media.SideCarHash = DataApi.GetStringP(res, 0, 35)
	token.Media.Blurhash = DataApi.GetStringP(res, 0, 36)
	return &token, nil
}

//修改完，未测试
func (r *queryResolver) ShareTokenValidatePassword(ctx context.Context, credentials models.ShareTokenCredentials) (bool, error) {
	var token models.ShareToken
	//if err := r.DB(ctx).Where("value = ?", credentials.Token).First(&token).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return false, errors.New("share not found")
	//	} else {
	//		return false, errors.Wrap(err, "failed to get share token from database")
	//	}
	//}

	sql_share_tokens_select := fmt.Sprintf("select * from share_tokens where value=%v limit 1", credentials.Token)
	dataApi, _ := DataApi.NewDataApiClient()
	res, _ := dataApi.Query(sql_share_tokens_select)
	if len(res) == 0 {
		return false, errors.New("share not found")
	}
	token.ID = DataApi.GetInt(res, 0, 0)
	token.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
	token.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 2)/1000, 0)
	token.Value = DataApi.GetString(res, 0, 3)
	token.OwnerID = DataApi.GetInt(res, 0, 4)
	expire := time.Unix(DataApi.GetLong(res, 0, 5)/1000, 0)
	token.Expire = &expire
	token.Password = DataApi.GetStringP(res, 0, 6)
	token.AlbumID = DataApi.GetIntP(res, 0, 7)
	token.MediaID = DataApi.GetIntP(res, 0, 8)
	if token.Password == nil {
		return true, nil
	}

	if credentials.Password == nil {
		return false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*token.Password), []byte(*credentials.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		} else {
			return false, errors.Wrap(err, "could not compare token password hashes")
		}
	}

	return true, nil
}

func (r *mutationResolver) ShareAlbum(ctx context.Context, albumID int, expire *time.Time, password *string) (*models.ShareToken, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return actions.AddAlbumShare(user, albumID, expire, password)
}

func (r *mutationResolver) ShareMedia(ctx context.Context, mediaID int, expire *time.Time, password *string) (*models.ShareToken, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return actions.AddMediaShare(user, mediaID, expire, password)
}

func (r *mutationResolver) DeleteShareToken(ctx context.Context, tokenValue string) (*models.ShareToken, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return actions.DeleteShareToken(user.ID, tokenValue)
}

func (r *mutationResolver) ProtectShareToken(ctx context.Context, tokenValue string, password *string) (*models.ShareToken, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return actions.ProtectShareToken(user.ID, tokenValue, password)
}
