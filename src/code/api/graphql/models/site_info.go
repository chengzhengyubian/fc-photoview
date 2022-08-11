package models

//修改完
import (
	"fmt"
	DataApi "github.com/photoview/photoview/api/dataapi"
	"github.com/pkg/errors"
)

type SiteInfo struct {
	InitialSetup         bool `gorm:"not null"`
	PeriodicScanInterval int  `gorm:"not null"`
	ConcurrentWorkers    int  `gorm:"not null"`
}

func (SiteInfo) TableName() string {
	return "site_info"
}

//	这里关注一下
func DefaultSiteInfo( /*db *gorm.DB*/ ) SiteInfo {
	defaultConcurrentWorkers := 3
	/*if db_drivers.SQLITE.MatchDatabase(db) {
		defaultConcurrentWorkers = 1
	}*/

	return SiteInfo{
		InitialSetup:         true,
		PeriodicScanInterval: 0,
		ConcurrentWorkers:    defaultConcurrentWorkers,
	}
}

// GetSiteInfo gets the site info row from the database, and creates it if it does not exist
//修改中，还剩下一部分没测试
func GetSiteInfo( /*db *gorm.DB*/ ) (*SiteInfo, error) {

	var siteInfo []*SiteInfo

	//if err := db.Limit(1).Find(&siteInfo).Error; err != nil { //SELECT * FROM `site_info` LIMIT 1
	//
	//	return nil, errors.Wrap(err, "get site info from database")
	//}

	sql_site_info_se := "SELECT * FROM `site_info` LIMIT 1"
	dataAPi, _ := DataApi.NewDataApiClient()
	res, err := dataAPi.Query(sql_site_info_se)
	if res == nil {
		return nil, errors.Wrap(err, "get site info from database")
	}
	num := len(res)
	for i := 0; i < num; i++ {
		var siteinfo SiteInfo
		siteinfo.InitialSetup = *res[i][0].BooleanValue
		siteinfo.PeriodicScanInterval = int(*res[i][1].LongValue)
		siteinfo.ConcurrentWorkers = int(*res[i][2].LongValue)
		siteInfo = append(siteInfo, &siteinfo)
	}
	if len(siteInfo) == 0 {
		newSiteInfo := DefaultSiteInfo( /*db*/ ) //初始化一张表,这里还没改
		//var setup int
		//if newSiteInfo.InitialSetup == true {
		//	setup = 1
		//} else {
		//	setup = 0
		//}
		//if err := db.Create(&newSiteInfo).Error; err != nil {
		//	return nil, errors.Wrap(err, "initialize site_info")
		//}
		//return &newSiteInfo, nil
		//sql_site_info_in := "INSERT INTO  site_info( initial_setup , periodic_scan_interval,concurrent_workers ) VALUES (" + strconv.Itoa(setup) + "," + strconv.Itoa(newSiteInfo.PeriodicScanInterval) + "," + strconv.Itoa(newSiteInfo.ConcurrentWorkers) + ")"
		sql_site_info_in := fmt.Sprintf("INSERT INTO  site_info( initial_setup , periodic_scan_interval,concurrent_workers ) VALUES(%v,%v,%v)", 1, newSiteInfo.PeriodicScanInterval, newSiteInfo.ConcurrentWorkers)
		dataAPi.ExecuteSQl(sql_site_info_in)
		return &newSiteInfo, nil
	} else {
		return siteInfo[0], nil
	}
	//if len(siteInfo) == 0 {
	//	newSiteInfo := DefaultSiteInfo(db)
	//
	//	if err := db.Create(&newSiteInfo).Error; err != nil {
	//		return nil, errors.Wrap(err, "initialize site_info")
	//	}
	//
	//	return &newSiteInfo, nil
	//} else {
	//	return siteInfo[0], nil
	//}

}
