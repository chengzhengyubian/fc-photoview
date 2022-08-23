package actions

import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strconv"

	//"gorm.io/gorm"
	"time"

	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

func AddMediaShare(user *models.User, mediaID int, expire *time.Time, password *string) (*models.ShareToken, error) {
	var media models.Media
	sql_media_select := fmt.Sprintf("select media.* from media left join albums on media.album_id=albums.id where EXISTS (SELECT * FROM user_albums WHERE user_albums.album_id = albums.id AND user_albums.user_id = %v) and media.id=%v limit 1", user.ID, mediaID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_media_select)
	if len(res) == 0 {
		return nil, err
	}
	//num:=len(res)
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

	if err != nil {
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//	return nil, auth.ErrUnauthorized
		//} else {
		return nil, errors.Wrap(err, "failed to validate media owner with database")
		//}
	}
	hashedPassword, err := hashSharePassword(password)
	if err != nil {
		return nil, err
	}

	shareToken := models.ShareToken{
		Value:    utils.GenerateToken(),
		OwnerID:  user.ID,
		Expire:   expire,
		Password: hashedPassword,
		AlbumID:  nil,
		MediaID:  &mediaID,
	}

	//timestr := shareToken.Expire.Format("2006-01-02 15:04:05")
	var timestr string
	if shareToken.Expire != nil {
		timestr = shareToken.Expire.Format("2006-01-02 15:04:05")
	} else {
		timestr = "NULL"
	}
	var pass string
	if shareToken.Password == nil {
		pass = "NULL"
	} else {
		pass = *shareToken.Password
	}
	var id string
	if shareToken.MediaID != nil {
		id = strconv.Itoa(*shareToken.MediaID)
	} else {
		id = "NULL"
	}
	sql_share_tokens_insert := fmt.Sprintf("insert into share_tokens (created_at, updated_at,value,owner_id,expire,password,album_id,media_id) values(NOW(),NOW(),'%v',%v,%v,'%v',NULL,%v)", shareToken.Value, shareToken.OwnerID, timestr, pass, id)
	dataApi.ExecuteSQl(sql_share_tokens_insert)
	return &shareToken, nil
}

func AddAlbumShare(user *models.User, albumID int, expire *time.Time, password *string) (*models.ShareToken, error) {
	var count int64

	sql_count_select := fmt.Sprintf("select count(*) from albums where EXISTS (SELECT * FROM user_albums WHERE user_albums.album_id = albums.id AND user_albums.user_id = %v)", user.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_count_select)
	count = DataApi.GetLong(res, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate album owner with database")
	}

	if count == 0 {
		return nil, auth.ErrUnauthorized
	}

	var hashedPassword *string = nil
	if password != nil {
		hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
		if err != nil {
			return nil, errors.Wrap(err, "failed to hash token password")
		}
		hashedStr := string(hashedPassBytes)
		hashedPassword = &hashedStr
	}

	shareToken := models.ShareToken{
		Value:    utils.GenerateToken(),
		OwnerID:  user.ID,
		Expire:   expire,
		Password: hashedPassword,
		AlbumID:  &albumID,
		MediaID:  nil,
	}

	var timestr string
	if shareToken.Expire != nil {
		timestr = shareToken.Expire.Format("2006-01-02 15:04:05")
	} else {
		timestr = "NULL"
	}
	var pass string
	if shareToken.Password == nil {
		pass = "NULL"
	} else {
		pass = *shareToken.Password
	}
	var id string
	if shareToken.AlbumID != nil {
		id = strconv.Itoa(*shareToken.AlbumID)
	} else {
		id = "NULL"
	}
	sql_share_tokens_insert := fmt.Sprintf("insert into share_tokens (created_at, updated_at,value,owner_id,expire,password,album_id,media_id) values(NOW(),NOW(),'%v',%v,%v,'%v',%v,NULL)", shareToken.Value, shareToken.OwnerID, timestr, pass, id)
	dataApi.ExecuteSQl(sql_share_tokens_insert)
	return &shareToken, nil
}

func DeleteShareToken(userID int, tokenValue string) (*models.ShareToken, error) {
	token, err := getUserToken(userID, tokenValue)
	if err != nil {
		return nil, err
	}

	sql_share_tokens_delete := fmt.Sprintf("delete from share_tokens where id=%v", token.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_share_tokens_delete)
	return token, nil
}

func ProtectShareToken(userID int, tokenValue string, password *string) (*models.ShareToken, error) {
	token, err := getUserToken(userID, tokenValue)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashSharePassword(password)
	if err != nil {
		return nil, err
	}

	token.Password = hashedPassword
	var pass string
	if token.Password == nil {
		pass = "NULL"
	} else {
		pass = *token.Password
	}

	sql_share_tokens_up := fmt.Sprintf("update share_tokens set password='%v", pass)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_share_tokens_up)
	return token, nil
}

func hashSharePassword(password *string) (*string, error) {
	var hashedPassword *string = nil
	if password != nil {
		hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate hash for share password")
		}
		hashedStr := string(hashedPassBytes)
		hashedPassword = &hashedStr
	}

	return hashedPassword, nil
}

func getUserToken(userID int, tokenValue string) (*models.ShareToken, error) {

	var token models.ShareToken

	sql_share_tokens_select := fmt.Sprintf("select share_tokens.* from share_tokens left join users on share_tokens.owner_id=users.id where share_tokens.value= '%v' and users.id = %v OR users.admin = TRUE limit 1", tokenValue, userID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_share_tokens_select)
	if len(res) == 0 {
		return nil, errors.Wrap(err, "failed to get user share token from database")
	}
	token.ID = DataApi.GetInt(res, 0, 0)
	token.CreatedAt = time.Unix(DataApi.GetLong(res, 0, 1)/1000, 0)
	token.UpdatedAt = time.Unix(DataApi.GetLong(res, 0, 2)/1000, 1)
	token.Value = DataApi.GetString(res, 0, 3)
	token.OwnerID = DataApi.GetInt(res, 0, 4)
	expire := time.Unix(DataApi.GetLong(res, 0, 5)/1000, 0)
	token.Expire = &expire
	token.Password = DataApi.GetStringP(res, 0, 6)
	token.AlbumID = DataApi.GetIntP(res, 0, 7)
	token.MediaID = DataApi.GetIntP(res, 0, 8)
	return &token, nil
}
