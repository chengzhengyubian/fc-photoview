package models

import "time"

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"rds-data-20220330/client"
	"strconv"
	"strings"
)

type User struct {
	Model
	Username string  `gorm:"unique;size:128"`
	Password *string `gorm:"size:256"`
	// RootPath string  `gorm:"size:512`
	Albums []Album `gorm:"many2many:user_albums;constraint:OnDelete:CASCADE;"`
	Admin  bool    `gorm:"default:false"`
}

type UserMediaData struct {
	ModelTimestamps
	UserID   int  `gorm:"primaryKey;autoIncrement:false"`
	MediaID  int  `gorm:"primaryKey;autoIncrement:false"`
	Favorite bool `gorm:"not null;default:false"`
}

type UserAlbums struct {
	UserID  int `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE;"`
	AlbumID int `gorm:"primaryKey;autoIncrement:false;constraint:OnDelete:CASCADE;"`
}

type AccessToken struct {
	Model
	UserID int       `gorm:"not null;index"`
	User   User      `gorm:"constraint:OnDelete:CASCADE;"`
	Value  string    `gorm:"not null;size:24;index"`
	Expire time.Time `gorm:"not null;index"`
}

type UserPreferences struct {
	Model
	UserID   int  `gorm:"not null;index"`
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
	Language *LanguageTranslation
}

func (u *UserPreferences) BeforeSave( /*tx *gorm.DB*/ ) error {

	if u.Language != nil && *u.Language == "" {
		u.Language = nil
	}

	if u.Language != nil {
		lang_str := string(*u.Language)
		found_match := false
		for _, lang := range AllLanguageTranslation {
			if string(lang) == lang_str {
				found_match = true
				break
			}
		}

		if !found_match {
			return errors.New("invalid language value")
		}
	}

	return nil
}

var ErrorInvalidUserCredentials = errors.New("invalid credentials")

//修改后
func AuthorizeUser(username string, password string) (*User, error) {
	var user User

	sql_users_se := "select * from users where username=\"" + username + "\""
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_se)
	if len(res.Body.Data.Records) == 0 {
		return nil, errors.Wrap(err, "failed to get user by username when authorizing")
	}
	user.ID = int(*res.Body.Data.Records[0][0].LongValue)
	user.Password = res.Body.Data.Records[0][4].StringValue
	if user.Password == nil {
		return nil, errors.New("user does not have a password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrorInvalidUserCredentials
		} else {
			return nil, errors.Wrap(err, "compare user password hash")
		}
	}

	return &user, nil
}

//修改后
func RegisterUser(username string, password *string, admin bool) (*User, error) {
	var res *client.ExecuteStatementResponse
	user := User{
		Username: username,
		Admin:    admin,
	}
	if password != nil {
		hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
		if err != nil {
			return nil, errors.Wrap(err, "failed to hash password")
		}
		hashedPass := string(hashedPassBytes)

		user.Password = &hashedPass
	}
	var ad int
	if admin == true {
		ad = 1
	} else {
		ad = 0
	}
	var sql_users_in string
	var pass string
	if password != nil {
		pass = *user.Password
	}
	if password != nil {
		sql_users_in = fmt.Sprintf("insert into users (created_at,updated_at,username,password,admin) values(NOW(),NOW(),'%v','%v',%v)", username, pass, ad)
	} else {
		sql_users_in = fmt.Sprintf("insert into users (created_at,updated_at,username,admin) values(NOW(),NOW(),'%v',%v)", username, ad)
	}
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_in)
	if err != nil {
		fmt.Println(res)
		return nil, errors.Wrap(err, "insert new user with password into database")
	}
	sql_users_se := "select * from users where username=\"" + username + "\""
	res, err = dataApi.ExecuteSQl(sql_users_se)
	if err != nil {
		fmt.Println(err)
	}
	user.ID = int(*res.Body.Data.Records[0][0].LongValue)
	return &user, nil
}

func (user *User) GenerateAccessToken() (*AccessToken, error) {

	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not generate token: %s\n", err.Error()))
	}
	const CHARACTERS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for i, b := range bytes {
		bytes[i] = CHARACTERS[b%byte(len(CHARACTERS))]
	}

	token_value := string(bytes)
	expire := time.Now().Add(14 * 24 * time.Hour)
	timeStr := expire.Format("2006-01-02 15:04:05")
	token := AccessToken{
		UserID: user.ID,
		Value:  token_value,
		Expire: expire,
	}
	sql_access_tokens_in := "insert into access_tokens(created_at,updated_at,user_id,value,expire) values(NOW(),NOW()," + strconv.Itoa(user.ID) + ",'" + token_value + "','" + timeStr + "')"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_access_tokens_in)
	fmt.Println(res)
	if err != nil {
		return nil, errors.Wrap(err, "saving access token to database")
	}
	return &token, nil
}

// FillAlbums fill user.Albums with albums from database 修改完
func (user *User) FillAlbums() error {
	// Albums already present
	if len(user.Albums) > 0 {
		return nil
	}
	sql_users_albums_se := "SELECT `albums`.`id`,`albums`.`created_at`,`albums`.`updated_at`,`albums`.`title`,`albums`.`parent_album_id`,`albums`.`path`,`albums`.`path_hash`,`albums`.`cover_id` FROM `albums` JOIN `user_albums` ON `user_albums`.`album_id` = `albums`.`id` AND `user_albums`.`user_id` =  " + strconv.Itoa(user.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_albums_se)
	fmt.Println(res)
	num := len(res.Body.Data.Records)
	for i := 0; i < num; i++ {
		var album Album
		album.ID = int(*res.Body.Data.Records[i][0].LongValue)
		album.CreatedAt = time.Unix(*res.Body.Data.Records[i][1].LongValue/1000, 0)
		album.UpdatedAt = time.Unix(*res.Body.Data.Records[i][2].LongValue/1000, 0)
		album.Title = *res.Body.Data.Records[i][3].StringValue
		if res.Body.Data.Records[i][4].IsNull != nil {
			album.ParentAlbumID = nil
		} else {
			id := int(*res.Body.Data.Records[i][4].LongValue)
			album.ParentAlbumID = &id
		}
		album.Path = *res.Body.Data.Records[i][5].StringValue
		album.PathHash = *res.Body.Data.Records[i][6].StringValue
		if res.Body.Data.Records[i][7].IsNull != nil {
			album.CoverID = nil
		} else {
			id := int(*res.Body.Data.Records[i][7].LongValue)
			album.CoverID = &id
		}
		user.Albums = append(user.Albums, album)
	}
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (user *User) OwnsAlbum(album *Album) (bool, error) {

	if err := user.FillAlbums(); err != nil {
		return false, err
	}

	albumIDs := make([]int, 0)
	for _, a := range user.Albums {
		albumIDs = append(albumIDs, a.ID)
	}
	albumID, _ := json.Marshal(albumIDs)
	albumids := strings.Trim(string(albumID), "[]")
	filter := func(sql string) string {
		return sql + fmt.Sprintf(" where id in (%v)", albumids)
	}
	ownedParents, err := album.GetParents(filter)
	if err != nil {
		return false, err
	}

	return len(ownedParents) > 0, nil
}

// FavoriteMedia sets/clears a media as favorite for the user
func (user *User) FavoriteMedia( /*db *gorm.DB,*/ mediaID int, favorite bool) (*Media, error) {
	var fav int
	if favorite == true {
		fav = 1
	} else {
		fav = 0
	}
	sql_user_media_data_se := fmt.Sprintf("select * from user_media_data where user_id =%v and media_id=%v", user.ID, mediaID)
	sql_user_media_data_in := "INSERT INTO `user_media_data` (`created_at`,`updated_at`,`user_id`,`media_id`,`favorite`) VALUES (NOW(),NOW()," + strconv.Itoa(user.ID) + "," + strconv.Itoa(mediaID) + "," + strconv.Itoa(fav) + ")"
	sql_user_media_data_up := fmt.Sprintf("UPDATE user_media_data set updated_at=NOW(),favorite=%v where user_id=%v and media_id=%v", fav, user.ID, mediaID)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_user_media_data_se)
	if len(res.Body.Data.Records) == 0 {
		res, err = dataApi.ExecuteSQl(sql_user_media_data_in)
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		res, err = dataApi.ExecuteSQl(sql_user_media_data_up)
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	}
	var media Media
	sql_media_se := "SELECT * FROM `media` WHERE `media`.`id` = " + strconv.Itoa(mediaID) //+ " ORDER BY `" + strconv.Itoa(mediaID) + "` LIMIT 1"
	res, err = dataApi.ExecuteSQl(sql_media_se)
	if err != nil {
		fmt.Println(nil)
	}
	fmt.Println(res)
	media.CreatedAt = time.Unix(*res.Body.Data.Records[0][1].LongValue/1000, 0)
	media.UpdatedAt = time.Unix(*res.Body.Data.Records[0][2].LongValue/1000, 0)
	media.Title = *res.Body.Data.Records[0][3].StringValue
	media.Path = *res.Body.Data.Records[0][4].StringValue
	media.PathHash = *res.Body.Data.Records[0][5].StringValue
	media.AlbumID = int(*res.Body.Data.Records[0][6].LongValue)
	if res.Body.Data.Records[0][7].IsNull != nil {
		media.ExifID = nil
	} else {
		id := int(*res.Body.Data.Records[0][7].LongValue)
		media.ExifID = &id
	}
	media.DateShot = time.Unix(*res.Body.Data.Records[0][8].LongValue/1000, 0)
	if *res.Body.Data.Records[0][9].StringValue == "photo" {
		media.Type = MediaTypePhoto
	} else {
		media.Type = MediaTypeVideo
	}
	if res.Body.Data.Records[0][10].IsNull != nil {
		media.VideoMetadataID = nil
	} else {
		id := int(*res.Body.Data.Records[0][10].LongValue)
		media.VideoMetadataID = &id
	}
	if res.Body.Data.Records[0][11].IsNull != nil {
		media.SideCarPath = nil
	} else {
		media.SideCarPath = res.Body.Data.Records[0][11].StringValue
	}
	if res.Body.Data.Records[0][12].IsNull != nil {
		media.SideCarHash = nil
	} else {
		media.SideCarHash = res.Body.Data.Records[0][12].StringValue
	}
	if res.Body.Data.Records[0][13].IsNull != nil {
		media.Blurhash = nil
	} else {
		media.Blurhash = res.Body.Data.Records[0][13].StringValue
	}
	return &media, nil
}
