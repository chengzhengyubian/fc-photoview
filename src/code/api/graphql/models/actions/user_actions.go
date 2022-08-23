package actions

import (
	"errors"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/utils"
)

func DeleteUser(userID int) (*models.User, error) {

	// make sure the last admin user is not deleted
	var adminUsers []*models.User

	sql_users_se := "select * from users where admin = true limit 2"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.ExecuteSQl(sql_users_se)
	num := len(res.Body.Data.Records)
	for i := 0; i < num; i++ {
		var user models.User
		user.ID = int(*res.Body.Data.Records[0][0].LongValue)
		user.Username = *res.Body.Data.Records[0][3].StringValue
		user.Password = res.Body.Data.Records[0][4].StringValue
		user.Admin = *res.Body.Data.Records[0][5].BooleanValue
		adminUsers = append(adminUsers, &user)
	}
	if err != nil {
		fmt.Println(err)
	}
	if len(adminUsers) == 1 && adminUsers[0].ID == userID {
		return nil, errors.New("deleting sole admin user is not allowed")
	}
	var user models.User
	deletedAlbumIDs := make([]int, 0)
	{
		sql_users_se = "SELECT * FROM `users` WHERE `users`.`id` =" + strconv.Itoa(userID) + " ORDER BY `users`.`id` LIMIT 1"
		res, err := dataApi.Query(sql_users_se)
		user.ID = int(*res[0][0].LongValue)
		user.Username = *res[0][3].StringValue
		user.Password = res[0][4].StringValue
		user.Admin = *res[0][5].BooleanValue
		userAlbums := user.Albums
		if err != nil {
			fmt.Println(err)
		}
		sql_albums_se := "SELECT `albums`.`id`,`albums`.`created_at`,`albums`.`updated_at`,`albums`.`title`,`albums`.`parent_album_id`,`albums`.`path`,`albums`.`path_hash`,`albums`.`cover_id` FROM `albums` JOIN `user_albums` ON `user_albums`.`album_id` = `albums`.`id` AND `user_albums`.`user_id` = " + strconv.Itoa(userID)
		res, err = dataApi.Query(sql_albums_se)
		num = len(res)
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
			userAlbums = append(userAlbums, album)
		}
		sql_user_albums_de := "DELETE FROM `user_albums` WHERE `user_albums`.`user_id` =" + strconv.Itoa(userID)
		dataApi.ExecuteSQl(sql_user_albums_de)
		for _, album := range userAlbums {
			var associatedUsers int
			sql_users_count_se := "SELECT count(*) FROM `users` JOIN `user_albums` ON `user_albums`.`user_id` = `users`.`id` AND `user_albums`.`album_id` =" + strconv.Itoa(album.ID)
			res, err = dataApi.Query(sql_users_count_se)
			associatedUsers = int(*res[0][0].LongValue)
			if associatedUsers == 0 {
				deletedAlbumIDs = append(deletedAlbumIDs, album.ID)
				sql_albums_de := "DELETE FROM `users` WHERE `users`.`id` =" + strconv.Itoa(userID)
				dataApi.ExecuteSQl(sql_albums_de)
			}
		}
		sql_users_de := "DELETE FROM `users` WHERE `users`.`id` =" + strconv.Itoa(userID)
		dataApi.ExecuteSQl(sql_users_de)
	}
	if err != nil {
		return nil, err
	}

	// If there is only one associated user, clean up the cache folder and delete the album row
	for _, deletedAlbumID := range deletedAlbumIDs {
		cachePath := path.Join(utils.MediaCachePath(), strconv.Itoa(int(deletedAlbumID)))
		if err := os.RemoveAll(cachePath); err != nil {
			return &user, err
		}
	}
	return &user, nil
}
