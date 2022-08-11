package resolvers

//修改完
import (
	"context"
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"strconv"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/periodic_scanner"
	"github.com/photoview/photoview/api/scanner/scanner_queue"
	"github.com/pkg/errors"
)

func (r *mutationResolver) ScanAll(ctx context.Context) (*models.ScannerResult, error) {
	err := scanner_queue.AddAllToQueue()
	if err != nil {
		return nil, err
	}

	startMessage := "Scanner started"

	return &models.ScannerResult{
		Finished: false,
		Success:  true,
		Message:  &startMessage,
	}, nil
}

//修改完，未测试
func (r *mutationResolver) ScanUser(ctx context.Context, userID int) (*models.ScannerResult, error) {
	var user models.User
	//if err := r.DB(ctx).First(&user, userID).Error; err != nil { //SELECT * FROM `users` WHERE `users`.`id` = 2 ORDER BY `users`.`id` LIMIT 1
	//	return nil, errors.Wrap(err, "get user from database")
	//}
	sql_users_se := "SELECT * FROM `users` WHERE `users`.`id` = " + strconv.Itoa(userID) + " ORDER BY `users`.`id` LIMIT 1"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_users_se)
	if len(res) == 0 {
		return nil, errors.Wrap(err, "get user from database")
	}
	user.ID = int(*res[0][0].LongValue)
	//user.CreatedAt=time.Unix(DataApi.GetLong(res,0,1)/1000,0)
	//user.UpdatedAt=time.Unix(DataApi.GetLong(res,0,2)/1000,0)
	user.Username = *res[0][3].StringValue
	user.Password = res[0][4].StringValue
	user.Admin = *res[0][5].BooleanValue
	scanner_queue.AddUserToQueue(&user)

	startMessage := "Scanner started"
	return &models.ScannerResult{
		Finished: false,
		Success:  true,
		Message:  &startMessage,
	}, nil
}

//修改完，未测试
func (r *mutationResolver) SetPeriodicScanInterval(ctx context.Context, interval int) (int, error) {
	//db := r.DB(ctx)
	if interval < 0 {
		return 0, errors.New("interval must be 0 or above")
	}

	//全局更新
	//if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.SiteInfo{}).Update("periodic_scan_interval", interval).Error; err != nil {
	//	return 0, err
	//}
	//这里注意一下
	sql_site_info_up := fmt.Sprintf("update site_info set periodic_scan_interval=%v", interval)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_site_info_up)
	var siteInfo models.SiteInfo
	//if err := db.First(&siteInfo).Error; err != nil {
	//	return 0, err
	//}
	sql_site_info_se := fmt.Sprintf("select * from site_info limit 1")
	res, err := dataApi.Query(sql_site_info_se)
	if len(res) == 0 {
		return 0, err
	}
	siteInfo.InitialSetup = DataApi.GetBoolean(res, 0, 0)
	siteInfo.PeriodicScanInterval = DataApi.GetInt(res, 0, 1)
	siteInfo.ConcurrentWorkers = DataApi.GetInt(res, 0, 2)
	periodic_scanner.ChangePeriodicScanInterval(time.Duration(siteInfo.PeriodicScanInterval) * time.Second)

	return siteInfo.PeriodicScanInterval, nil
}

//修改完，未测试
func (r *mutationResolver) SetScannerConcurrentWorkers(ctx context.Context, workers int) (int, error) {
	//db := r.DB(ctx)
	//if workers < 1 {
	//return 0, errors.New("concurrent workers must at least be 1")
	//}

	//注意一下
	//if workers > 1 && drivers.DatabaseDriverFromEnv() == drivers.SQLITE {
	//	return 0, errors.New("multiple workers not supported for SQLite databases")
	//}
	//if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.SiteInfo{}).Update("concurrent_workers", workers).Error; err != nil {
	//	return 0, err
	//}
	sql_site_info_up := fmt.Sprintf("update site_info set concurrent_workers=%v", workers)
	dataApi, _ := DataApi.NewDataApiClient()
	dataApi.ExecuteSQl(sql_site_info_up)

	var siteInfo models.SiteInfo
	//if err := db.First(&siteInfo).Error; err != nil {
	//	return 0, err
	//}
	sql_site_info_se := fmt.Sprintf("select * from site_info limit 1")
	res, err := dataApi.Query(sql_site_info_se)
	if len(res) == 0 {
		return 0, err
	}
	siteInfo.InitialSetup = DataApi.GetBoolean(res, 0, 0)
	siteInfo.PeriodicScanInterval = DataApi.GetInt(res, 0, 1)
	siteInfo.ConcurrentWorkers = DataApi.GetInt(res, 0, 2)

	scanner_queue.ChangeScannerConcurrentWorkers(siteInfo.ConcurrentWorkers)

	//db := r.DB(ctx)
	//if workers < 1 {
	//	return 0, errors.New("concurrent workers must at least be 1")
	//}
	//
	//if workers > 1 && drivers.DatabaseDriverFromEnv() == drivers.SQLITE {
	//	return 0, errors.New("multiple workers not supported for SQLite databases")
	//}
	//
	//if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&models.SiteInfo{}).Update("concurrent_workers", workers).Error; err != nil {
	//	return 0, err
	//}
	//
	//var siteInfo models.SiteInfo
	//if err := db.First(&siteInfo).Error; err != nil {
	//	return 0, err
	//}
	//
	//scanner_queue.ChangeScannerConcurrentWorkers(siteInfo.ConcurrentWorkers)

	return siteInfo.ConcurrentWorkers, nil
}
