package resolvers

//修改完
import (
	"context"
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	_ "github.com/photoview/photoview/api/database"
	api "github.com/photoview/photoview/api/graphql"
	"github.com/photoview/photoview/api/graphql/auth"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/graphql/models/actions"
	"github.com/photoview/photoview/api/scanner"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/photoview/photoview/api/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
	"log"
	"os"
	"path"
	"rds-data-20220330/client"
	"strconv"
	"strings"
	"time"
)

type userResolver struct {
	*Resolver
}

func (r *Resolver) User() api.UserResolver {
	return &userResolver{r}
}

//修改完
func (r *queryResolver) User(ctx context.Context, order *models.Ordering, paginate *models.Pagination) ([]*models.User, error) {
	sql_users_se := "select * from users"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_se)
	if err != nil {
		fmt.Println(err)
	}
	num1 := len(res.Body.Data.Records)
	var users []*models.User
	for i := 0; i < num1; i++ {
		var user models.User
		user.ID = int(*res.Body.Data.Records[i][0].LongValue)
		user.Password = res.Body.Data.Records[i][4].StringValue
		user.Username = *res.Body.Data.Records[i][3].StringValue
		user.Admin = *res.Body.Data.Records[i][5].BooleanValue
		users = append(users, &user)
	}
	return users, nil
}

//
func (r *userResolver) Albums(ctx context.Context, user *models.User) ([]*models.Album, error) {
	user.FillAlbums()

	pointerAlbums := make([]*models.Album, len(user.Albums))
	for i, album := range user.Albums {
		pointerAlbums[i] = &album
	}
	return pointerAlbums, nil
}

//修改完
func (r *userResolver) RootAlbums(ctx context.Context, user *models.User) (albums []*models.Album, err error) {
	//SELECT `albums`.`id`,`albums`.`created_at`,`albums`.`updated_at`,`albums`.`title`,`albums`.`parent_album_id`,`albums`.`path`,`albums`.`path_hash`,`albums`.`cover_id` FROM `albums` JOIN `user_albums` ON `user_albums`.`album_id` = `albums`.`id` AND `user_albums`.`user_id` = 5 WHERE albums.parent_album_id NOT IN (SELECT albums.id FROM `user_albums` JOIN albums ON albums.id = user_albums.album_id AND user_albums.user_id = 5) OR albums.parent_album_id IS NULL
	fmt.Println(user.ID)
	dataApi, _ := DataApi.NewDataApiClient()
	sql := "SELECT `albums`.`id`,`albums`.`created_at`,`albums`.`updated_at`,`albums`.`title`,`albums`.`parent_album_id`,`albums`.`path`,`albums`.`path_hash`,`albums`.`cover_id` FROM `albums` JOIN `user_albums` ON `user_albums`.`album_id` = `albums`.`id` AND `user_albums`.`user_id` =" + strconv.Itoa(user.ID) + " WHERE albums.parent_album_id NOT IN (SELECT albums.id FROM `user_albums` JOIN albums ON albums.id = user_albums.album_id AND user_albums.user_id = " + strconv.Itoa(user.ID) + ") OR albums.parent_album_id IS NULL"
	res, err := dataApi.Query(sql)
	fmt.Println(res)
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
	fmt.Println(len(albums))
	return
}

func (r *queryResolver) MyUser(ctx context.Context) (*models.User, error) {

	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}

	return user, nil
}

//修改完
func (r *mutationResolver) AuthorizeUser(ctx context.Context, username string, password string) (*models.AuthorizeResult, error) {
	//db := r.DB(ctx)
	user, err := models.AuthorizeUser(username, password)
	if err != nil {
		return &models.AuthorizeResult{
			Success: false,
			Status:  err.Error(),
		}, nil
	}
	var token *models.AccessToken
	token, err = user.GenerateAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	return &models.AuthorizeResult{
		Success: true,
		Status:  "ok",
		Token:   &token.Value,
	}, nil
}

//修改中
func (r *mutationResolver) InitialSetupWizard(ctx context.Context, username string, password string, rootPath string) (*models.AuthorizeResult, error) {
	//db := r.DB(ctx)
	siteInfo, err := models.GetSiteInfo( /*db*/ )
	if err != nil {
		return nil, err
	}

	if !siteInfo.InitialSetup {
		return nil, errors.New("not initial setup")
	}

	rootPath = path.Clean(rootPath)

	var token *models.AccessToken
	{
		//transactionError := db.Transaction(func(tx *gorm.DB) error {
		//if err := tx.Exec("UPDATE site_info SET initial_setup = false").Error; err != nil {
		//	return err
		//}
		sql_site_info_up := "UPDATE site_info SET initial_setup = false"
		dataApi, _ := DataApi.NewDataApiClient()
		res, err := dataApi.ExecuteSQl(sql_site_info_up)
		if err != nil {
			fmt.Println(res)
			fmt.Println(err)
		}
		user, err := models.RegisterUser(username, &password, true)
		if err != nil {
			return nil, err
		}

		_, err = scanner.NewRootAlbum(rootPath, user)
		if err != nil {
			return nil, err
		}

		token, err = user.GenerateAccessToken()
		if err != nil {
			return nil, err
		}

		//	return nil
		//	})
	}
	//if transactionError != nil {
	//	return &models.AuthorizeResult{
	//		Success: false,
	//		Status:  err.Error(),
	//	}, nil
	//}

	return &models.AuthorizeResult{
		Success: true,
		Status:  "ok",
		Token:   &token.Value,
	}, nil
}

//修改完
func (r *queryResolver) MyUserPreferences(ctx context.Context) (*models.UserPreferences, error) {
	dataApi, _ := DataApi.NewDataApiClient()
	var res *client.ExecuteStatementResponse
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}
	userPref := models.UserPreferences{
		UserID: user.ID,
	}
	id := strconv.Itoa(user.ID)
	sql1 := "select * from user_preferences where user_id=" + id
	res, err := dataApi.ExecuteSQl(sql1)
	if err != nil {
		fmt.Println(err)
	}
	if len(res.Body.Data.Records) == 0 {
		sql2 := "insert into user_preferences (user_id,language,created_at) VALUES (" + id + ",\"English\",NOW())"
		dataApi.ExecuteSQl(sql2)
		sql := "select * from user_preferences where user_id=" + id
		res, err := dataApi.ExecuteSQl(sql)
		if err != nil {
			fmt.Println(err)
		}
		var language *string
		var langTrans *models.LanguageTranslation = nil
		language = res.Body.Data.Records[0][4].StringValue
		lng := models.LanguageTranslation(*language)
		langTrans = &lng
		userPref.Language = langTrans
		userPref.ID = int(*res.Body.Data.Records[0][0].LongValue)
	} else {
		var language *string
		language = res.Body.Data.Records[0][4].StringValue
		var langTrans *models.LanguageTranslation = nil
		if language != nil {
			lng := models.LanguageTranslation(*language)
			langTrans = &lng
		}
		userPref.Language = langTrans
		userPref.ID = int(*res.Body.Data.Records[0][0].LongValue)
	}
	return &userPref, nil
}

//修改完
func (r *mutationResolver) ChangeUserPreferences(ctx context.Context, language *string) (*models.UserPreferences, error) {
	dataApi, _ := DataApi.NewDataApiClient()
	var m *client.ExecuteStatementResponse
	var str string
	str = *language
	user := auth.UserFromContext(ctx)
	if user == nil {
		return nil, auth.ErrUnauthorized
	}
	var langTrans *models.LanguageTranslation = nil
	if language != nil {
		lng := models.LanguageTranslation(*language)
		langTrans = &lng
	}
	var userPref models.UserPreferences
	//下面是对这段逻辑的改写
	var id string
	id = strconv.Itoa(user.ID)
	sql2 := "select * from user_preferences where user_id =" + id
	//如果执行这个操作之后没有值，那么执行下面的操作，否则啥也不干
	m, err := dataApi.ExecuteSQl(sql2)
	if err != nil {
		fmt.Println(err)
	}
	if len(m.Body.Data.Records) == 0 {
		//执行，插入的操作
		sql3 := "insert into user_preferences (user_id,language,created_at) VALUES (" + id + ",\"English\",NOW())"
		dataApi.ExecuteSQl(sql3)
	}
	//更新的操作
	sql := "update user_preferences set (language,updated_at) values(\"" + str + "\",NOW()) where user_id=" + id
	dataApi.ExecuteSQl(sql)
	m, _ = dataApi.ExecuteSQl(sql2)
	userPref.ID = int(*m.Body.Data.Records[0][0].LongValue)
	userPref.Language = langTrans
	userPref.UserID = user.ID
	return &userPref, nil
}

//修改完
func (r *mutationResolver) UpdateUser(ctx context.Context, id int, username *string, password *string, admin *bool) (*models.User, error) {
	if username == nil && password == nil && admin == nil {
		return nil, errors.New("no updates requested")
	}

	sqlUsersInsert := DataApi.FormatSql("select * from users where id =%v", id)
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sqlUsersInsert)
	log.Print("insert user result: ", res)
	if err != nil {
		return nil, err
	}

	var user models.User
	user.Username = DataApi.GetString(res, 0, 3)
	user.Password = DataApi.GetStringP(res, 0, 4)
	user.Admin = DataApi.GetBoolean(res, 0, 5)
	if password != nil {
		hashedPassBytes, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
		if err != nil {
			return nil, err
		}
		hashedPass := string(hashedPassBytes)

		user.Password = &hashedPass
	}

	if username != nil {
		user.Username = *username
	}

	ad := 0
	if admin != nil {
		user.Admin = *admin
		if user.Admin {
			ad = 1
		}
	}

	updateData := make(map[string]any)
	if username != nil {
		updateData["username"] = username
	}
	if password != nil {
		updateData["password"] = user.Password
	}
	updateData["admin"] = ad

	updateWhere := make(map[string]any)
	updateWhere["id"] = id

	sqlUsersUpdate, err := DataApi.FormatUpdateSql("users", updateData, updateWhere)
	if err != nil {
		return nil, err
	}

	_, err = dataApi.ExecuteSQl(sqlUsersUpdate)
	log.Print("update user result: ", res)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, username string, password *string, admin bool) (*models.User, error) {
	var user *models.User
	var err error
	user, err = models.RegisterUser(username, password, admin)
	if err != nil {
		fmt.Println(err)
	}
	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*models.User, error) {
	return actions.DeleteUser( /*r.DB(ctx),*/ id)
}

//修改完
func (r *mutationResolver) UserAddRootPath(ctx context.Context, id int, rootPath string) (*models.Album, error) {
	//db := r.DB(ctx)

	rootPath = path.Clean(rootPath)

	var user models.User
	//if err := db.First(&user, id).Error; err != nil {
	//	return nil, err
	//}
	sql_users_se := "select * from users where id=" + strconv.Itoa(id) //+ "order by id limit 1"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_se)
	if len(res.Body.Data.Records) == 0 {
		return nil, err
	} else {
		user.ID = int(*res.Body.Data.Records[0][0].LongValue)
		user.Username = *res.Body.Data.Records[0][3].StringValue
		user.Password = res.Body.Data.Records[0][4].StringValue
		fmt.Println(res.Body.Data.Records[0][5].IsNull)
		user.Admin = *res.Body.Data.Records[0][5].BooleanValue
	}
	newAlbum, err := scanner.NewRootAlbum(rootPath, &user)
	if err != nil {
		return nil, err
	}

	return newAlbum, nil
}

//未改完有问题,去掉事务之后可以
func (r *mutationResolver) UserRemoveRootAlbum(ctx context.Context, userID int, albumID int) (*models.Album, error) {
	//db := r.DB(ctx)
	var album models.Album
	//if err := db.First(&album, albumID).Error; err != nil { //SELECT * FROM `albums` WHERE `albums`.`id` = 110 ORDER BY `albums`.`id` LIMIT 1
	//	return nil, err
	//}
	dataApi, _ := DataApi.NewDataApiClient()
	//sql_albums_se := "SELECT * FROM `albums` WHERE `albums`.`id`=" + strconv.Itoa(albumID) //+ "LIMIT 1"
	sql_albums_se := fmt.Sprintf("SELECT * FROM `albums` WHERE `albums`.`id`=%v limit 1", albumID)
	res, err := dataApi.Query(sql_albums_se)
	fmt.Println(res)
	if len(res) == 0 {
		return nil, errors.WithMessage(err, "can't find album")
	}
	var deletedAlbumIDs []int = nil
	//transactionError := db.Transaction(
	//	func(tx *gorm.DB) error {
	//if err := tx.Raw("DELETE FROM user_albums WHERE user_id = ? AND album_id = ?", userID, albumID).Error; err != nil {
	//	return err
	//}
	//sql_user_albums_de := "DELETE FROM user_albums WHERE user_id =" + strconv.Itoa(userID) + "AND album_id =" + strconv.Itoa(albumID)
	sql_user_albums_de := fmt.Sprintf("DELETE FROM user_albums WHERE user_id =%v and album_id =%v", userID, albumID)
	dataApi.ExecuteSQl(sql_user_albums_de)
	children, err := album.GetChildren(nil) //还未改
	if err != nil {
		return nil, err
	}
	childAlbumIDs := make([]int, len(children))
	for i, child := range children {
		childAlbumIDs[i] = child.ID
	}
	//result := tx.Exec("DELETE FROM user_albums WHERE user_id = ? and album_id IN (?)", userID, childAlbumIDs) //DELETE FROM user_albums WHERE user_id = 7 and album_id IN (110)
	childAlbumid, err := json.Marshal(childAlbumIDs)
	child := strings.Trim(string(childAlbumid), "[]")
	sql_user_albums_child_de := "DELETE FROM user_albums WHERE user_id =" + strconv.Itoa(userID) + " and album_id IN (" + child + ")"
	dataApi.ExecuteSQl(sql_user_albums_child_de)
	//if result.Error != nil { //SELECT * FROM `users` WHERE `users`.`id` = 13 ORDER BY `users`.`id` LIMIT 1
	//	return result.Error
	//}
	//if result.RowsAffected == 0 {
	//	return errors.New("No relation deleted")
	//}
	// Cleanup if no user owns the album anymore
	var userAlbumCount int
	//if err := tx.Raw("SELECT COUNT(user_id) FROM user_albums WHERE album_id = ?", albumID).Scan(&userAlbumCount).Error; err != nil { //SELECT COUNT(user_id) FROM user_albums WHERE album_id = 110
	//	return err
	//}
	sql_count_se := "SELECT COUNT(user_id) FROM user_albums WHERE album_id =" + strconv.Itoa(albumID)
	res, err = dataApi.Query(sql_count_se)
	if len(res) == 0 {
		return nil, err
	}
	userAlbumCount = int(*res[0][0].LongValue)
	if userAlbumCount == 0 {
		deletedAlbumIDs = append(childAlbumIDs, albumID)
		childAlbumIDs = nil

		//Delete albums from database
		//if err := tx.Delete(&models.Album{}, "id IN (?)", deletedAlbumIDs).Error; err != nil { // DELETE FROM `albums` WHERE id IN (110,110)
		//	return err
		//}
		del, err := json.Marshal(deletedAlbumIDs)
		if err != nil {
			return nil, err
		}
		de := strings.Trim(string(del), "[]")
		sql_albums_de := "DELETE FROM `albums` WHERE id IN (" + de + ")"
		dataApi.ExecuteSQl(sql_albums_de)
		deletedAlbumIDs = nil
	}
	//return nil
	//	//})
	//if transactionError != nil {
	//	return nil, transactionError
	//}
	if deletedAlbumIDs != nil {
		// Delete albums from cache
		for _, id := range deletedAlbumIDs {
			cacheAlbumPath := path.Join(utils.MediaCachePath(), strconv.Itoa(id))

			if err := os.RemoveAll(cacheAlbumPath); err != nil {
				return nil, err
			}
		}
		// Reload faces as media might have been deleted
		if face_detection.GlobalFaceDetector != nil {
			if err := face_detection.GlobalFaceDetector.ReloadFacesFromDatabase( /*db*/ ); err != nil {
				return nil, err
			}
		}
	}
	return &album, nil
}
