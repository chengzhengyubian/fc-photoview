package periodic_scanner

import (
	DataApi "github.com/photoview/photoview/api/dataapi"
	"log"
	"sync"
	"time"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/scanner_queue"
)

type periodicScanner struct {
	ticker         *time.Ticker
	ticker_changed chan bool
	mutex          *sync.Mutex
	/**/ /*db *gorm.DB*/
}

var mainPeriodicScanner *periodicScanner = nil

func getPeriodicScanInterval() (time.Duration, error) {

	var siteInfo models.SiteInfo
	//if err := db.First(&siteInfo).Error; err != nil { //SELECT * FROM `site_info` ORDER BY `site_info`.`initial_setup` LIMIT 1
	//	return 0, err
	//}
	sql_site_info_se := "SELECT * FROM `site_info` ORDER BY `site_info`.`initial_setup` LIMIT 1"
	dataApi, _ := DataApi.NewDataApiClient()
	res, err := dataApi.Query(sql_site_info_se)
	if len(res) == 0 {
		return 0, err
	}
	siteInfo.InitialSetup = *res[0][0].BooleanValue
	siteInfo.PeriodicScanInterval = int(*res[0][1].LongValue)
	siteInfo.ConcurrentWorkers = int(*res[0][1].LongValue)
	return time.Duration(siteInfo.PeriodicScanInterval) * time.Second, nil
}

func InitializePeriodicScanner( /*db *gorm.DB*/ ) error {
	if mainPeriodicScanner != nil {
		panic("periodic scanner has already been initialized")
	}

	scanInterval, err := getPeriodicScanInterval()
	if err != nil {
		return err
	}

	mainPeriodicScanner = &periodicScanner{
		//db:             db,
		ticker_changed: make(chan bool),
		mutex:          &sync.Mutex{},
	}

	go scanIntervalRunner()

	ChangePeriodicScanInterval(scanInterval)
	return nil
}

func ChangePeriodicScanInterval(duration time.Duration) {
	var new_ticker *time.Ticker = nil
	if duration > 0 {
		new_ticker = time.NewTicker(duration)
		log.Printf("Periodic scan interval changed: %s", duration.String())
	} else {
		log.Print("Periodic scan interval changed: disabled")
	}

	{
		mainPeriodicScanner.mutex.Lock()
		defer mainPeriodicScanner.mutex.Unlock()

		if mainPeriodicScanner.ticker != nil {
			mainPeriodicScanner.ticker.Stop()
		}

		mainPeriodicScanner.ticker = new_ticker
		mainPeriodicScanner.ticker_changed <- true
	}
}

func scanIntervalRunner() {
	for {
		log.Print("Scan interval runner: Waiting for signal")
		if mainPeriodicScanner.ticker != nil {
			select {
			case <-mainPeriodicScanner.ticker_changed:
				log.Print("Scan interval runner: New ticker detected")
			case <-mainPeriodicScanner.ticker.C:
				log.Print("Scan interval runner: Starting periodic scan")
				scanner_queue.AddAllToQueue()
			}
		} else {
			<-mainPeriodicScanner.ticker_changed
			log.Print("Scan interval runner: New ticker detected")
		}
	}
}
