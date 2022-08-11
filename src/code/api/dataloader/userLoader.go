package dataloader

/*修改完*/
import (
	"encoding/json"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strings"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
)

//修改完
func NewUserLoaderByToken( /*db *gorm.DB*/ ) *UserLoader {
	return &UserLoader{
		maxBatch: 100,
		wait:     5 * time.Millisecond,
		fetch: func(tokens []string) ([]*models.User, []error) {

			var accessTokens []*models.AccessToken
			//SELECT * FROM `access_tokens` WHERE expire > '2022-07-28 00:24:13.606' AND value IN ('WohDNhR2teZ344ldk4jkJyQL')
			//err := db.Where("expire > ?", time.Now()).Where("value IN (?)", tokens).Find(&accessTokens).Error
			dataApi, _ := DataApi.NewDataApiClient()
			expire := time.Now()
			timeStr := expire.Format("2006-01-02 15:04:05")
			token, err := json.Marshal(tokens)
			Token := strings.Trim(string(token), "[]")
			sql_access_tokens_se := "SELECT * FROM `access_tokens` WHERE expire > \"" + timeStr + "\" AND value IN (" + Token + ")"
			res, err := dataApi.Query(sql_access_tokens_se)
			fmt.Println(res)
			num := len(res)
			for i := 0; i < num; i++ {
				var accessToken models.AccessToken
				accessToken.ID = int(*res[i][0].LongValue)
				accessToken.CreatedAt = time.Unix(*res[i][1].LongValue/1000, 0)
				accessToken.UpdatedAt = time.Unix(*res[i][2].LongValue/1000, 0)
				accessToken.UserID = int(*res[i][3].LongValue)
				accessToken.Value = *res[i][4].StringValue
				accessToken.Expire = time.Unix(*res[i][5].LongValue/1000, 0)
				accessTokens = append(accessTokens, &accessToken)
			}
			if err != nil {
				return nil, []error{err}
			}
			//SELECT distinct user_id FROM `access_tokens` WHERE expire > '2022-07-28 00:24:13.682' AND value IN ('WohDNhR2teZ344ldk4jkJyQL')
			//rows, err := db.Table("access_tokens").Select("distinct user_id").Where("expire > ?", time.Now()).Where("value IN (?)", tokens).Rows()
			sql_access_tokens_se = "SELECT distinct user_id FROM `access_tokens` WHERE expire > \"" + timeStr + "\" AND value IN (" + Token + ")"
			res, err = dataApi.Query(sql_access_tokens_se)
			if err != nil {
				return nil, []error{err}
			}
			userIDs := make([]int, 0)
			num = len(res)
			for i := 0; i < num; i++ {
				id := int(*res[i][0].LongValue)
				userIDs = append(userIDs, id)
			}
			//for rows.Next() {
			//	var id int
			//	if err := db.ScanRows(rows, &id); err != nil {
			//		return nil, []error{err}
			//	} //这里关注一下，看看怎么测试
			//	userIDs = append(userIDs, id)
			//}
			//rows.Close()

			var userMap map[int]*models.User
			if len(userIDs) > 0 {

				var users []*models.User
				//if err := db.Where("id IN (?)", userIDs).Find(&users).Error; err != nil { // SELECT * FROM `users` WHERE id IN (2)
				//	return nil, []error{err}
				//}
				userID, err := json.Marshal(userIDs)
				if err != nil {
					fmt.Print(err)
				}
				userids := strings.Trim(string(userID), "[]")
				sql_users_se := "SELECT * FROM `users` WHERE id IN (" + userids + ")"
				res, err = dataApi.Query(sql_users_se)
				num = len(res)
				for i := 0; i < num; i++ {
					var user models.User
					user.ID = int(*res[i][0].LongValue)
					user.Password = res[i][4].StringValue
					user.Username = *res[i][3].StringValue
					user.Admin = *res[i][5].BooleanValue
					users = append(users, &user)
				}
				userMap = make(map[int]*models.User, len(users))
				for _, user := range users {
					userMap[user.ID] = user
				}
			} else {
				userMap = make(map[int]*models.User, 0)
			}

			tokenMap := make(map[string]*models.AccessToken, len(tokens))
			for _, token := range accessTokens {
				tokenMap[token.Value] = token
			}

			result := make([]*models.User, len(tokens))
			for i, token := range tokens {
				accessToken, tokenFound := tokenMap[token]
				if tokenFound {
					user, userFound := userMap[accessToken.UserID]
					if userFound {
						result[i] = user
					}
				}
			}

			return result, nil
		},
	}
}
