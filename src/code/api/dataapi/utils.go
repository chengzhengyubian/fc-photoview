package dataapi

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"rds-data-20220330/client"
	"reflect"
	"strings"
	"sync"
)

func FormatSql(template string, args ...any) string {
	// todo SQL防注入
	return fmt.Sprintf(template, args...)
}

// 获取数据
func GetString(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) string {
	if recordsDataIsNull(records, row, colum) {
		return ""
	}
	return *records[row][colum].StringValue
}

func GetStringP(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) *string {
	if recordsDataIsNull(records, row, colum) {
		return nil
	}
	return records[row][colum].StringValue
}
func GetInt(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) int {
	if recordsDataIsNull(records, row, colum) {
		return 0
	}
	return int(*records[row][colum].LongValue)
}
func GetIntP(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) *int {
	var n *int
	if recordsDataIsNull(records, row, colum) {
		return nil
	}
	m := int(*records[row][colum].LongValue)
	n = &m
	return n
}

func GetLong(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) int64 {
	if recordsDataIsNull(records, row, colum) {
		return 0
	}
	return *records[row][colum].LongValue
}
func GetLongP(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) *int64 {
	var n *int64
	var m int64
	if recordsDataIsNull(records, row, colum) {
		return nil
	}
	m = *records[row][colum].LongValue
	n = &m
	return n
}

func GetBoolean(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) bool {
	if recordsDataIsNull(records, row, colum) {
		return false
	}
	return *records[row][colum].BooleanValue
}

func recordsDataIsNull(records [][]*client.ExecuteStatementResponseBodyDataRecords, row int, colum int) bool {
	if records[row][colum].IsNull == nil {
		return false
	}
	return *records[row][colum].IsNull
}

// FormatUpdateSql /**
func FormatUpdateSql(table string, data map[string]any, query map[string]any) (string, error) {
	if len(data) == 0 {
		return "", errors.New(fmt.Sprintf("table %s has no data to update", table))
	}

	update := ""
	for k, v := range data {
		elem, err := formatSqlElem(k, v)
		if err != nil {
			return "", err
		}
		update = update + elem + ", "
	}
	update = strings.TrimSuffix(update, ", ")

	where := "1=1 AND "
	for k, v := range query {
		elem, err := formatSqlElem(k, v)
		if err != nil {
			return "", err
		}
		where = where + elem + " AND "
	}
	where = strings.TrimSuffix(where, "AND ")

	sql := fmt.Sprintf("update %s set %s where %s", table, update, where)
	log.Printf("update sql is: %s", sql)

	return sql, nil
}

func formatSqlElem(key string, value any) (string, error) {
	switch value.(type) {
	case string:
		return FormatSql("%s='%s'", key, value), nil
	case *string:
		return FormatSql("%s='%s'", key, *value.(*string)), nil
	case int:
		return FormatSql("%s=%d", key, value), nil
	default:
		return "", errors.New(fmt.Sprintf("datatype %s not supported yet!", reflect.TypeOf(value)))
	}
}

func MigrateDatabase() {
	sql_site_info_insert := "CREATE TABLE if not exists `site_info` (\n`initial_setup` tinyint(1) NOT NULL,\n`periodic_scan_interval` bigint(20) NOT NULL,\n`concurrent_workers` bigint(20) NOT NULL\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_users_insert := "CREATE TABLE if not exists `users` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`username` varchar(128) DEFAULT NULL,\n`password` varchar(256) DEFAULT NULL,\n`admin` tinyint(1) DEFAULT '0',\nPRIMARY KEY (`id`),\nUNIQUE KEY `username` (`username`)\n) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_access_tokens_insert := "CREATE TABLE if not exists `access_tokens` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`user_id` bigint(20) NOT NULL,\n`value` varchar(24) NOT NULL,\n`expire` datetime(3) NOT NULL,\nPRIMARY KEY (`id`),\nKEY `idx_access_tokens_user_id` (`user_id`),\nKEY `idx_access_tokens_value` (`value`),\nKEY `idx_access_tokens_expire` (`expire`),\nCONSTRAINT `fk_access_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci "
	sql_user_preferences_insert := "CREATE TABLE if not exists `user_preferences` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`user_id` bigint(20) NOT NULL,\n`language` longtext,\nPRIMARY KEY (`id`),\nKEY `idx_user_preferences_user_id` (`user_id`),\nCONSTRAINT `fk_user_preferences_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_albums_insert := "CREATE TABLE if not exists `albums` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`title` longtext NOT NULL,\n`parent_album_id` bigint(20) DEFAULT NULL,\n`path` longtext NOT NULL,\n`path_hash` varchar(191) DEFAULT NULL,\n`cover_id` bigint(20) DEFAULT NULL,\nPRIMARY KEY (`id`),\nUNIQUE KEY `path_hash` (`path_hash`),\nKEY `idx_albums_parent_album_id` (`parent_album_id`),\nCONSTRAINT `fk_albums_parent_album` FOREIGN KEY (`parent_album_id`) REFERENCES `albums` (`id`) ON DELETE SET NULL\n) ENGINE=InnoDB AUTO_INCREMENT=315 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_user_albums_insert := "CREATE TABLE if not exists `user_albums` (\n`album_id` bigint(20) NOT NULL,\n`user_id` bigint(20) NOT NULL,\nPRIMARY KEY (`album_id`,`user_id`),\nKEY `fk_user_albums_user` (`user_id`),\nCONSTRAINT `fk_user_albums_album` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_user_albums_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_video_metadata_insert := "CREATE TABLE if not exists `video_metadata` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`width` bigint(20) NOT NULL,\n`height` bigint(20) NOT NULL,\n`duration` double NOT NULL,\n`codec` longtext,\n`framerate` double DEFAULT NULL,\n`bitrate` longtext,\n`color_profile` longtext,\n`audio` longtext,\nPRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_media_exif_insert := "CREATE TABLE if not exists `media_exif` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`camera` longtext,\n`maker` longtext,\n`lens` longtext,\n`date_shot` datetime(3) DEFAULT NULL,\n`exposure` double DEFAULT NULL,\n`aperture` double DEFAULT NULL,\n`iso` bigint(20) DEFAULT NULL,\n`focal_length` double DEFAULT NULL,\n`flash` bigint(20) DEFAULT NULL,\n`orientation` bigint(20) DEFAULT NULL,\n`exposure_program` bigint(20) DEFAULT NULL,\n`gps_latitude` double DEFAULT NULL,\n`gps_longitude` double DEFAULT NULL,\n`description` longtext,\nPRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_media_insert := "CREATE TABLE if not exists `media` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`title` longtext NOT NULL,\n`path` longtext NOT NULL,\n`path_hash` varchar(191) NOT NULL,\n`album_id` bigint(20) NOT NULL,\n`exif_id` bigint(20) DEFAULT NULL,\n`date_shot` datetime(3) NOT NULL,\n`type` varchar(191) NOT NULL,\n`video_metadata_id` bigint(20) DEFAULT NULL,\n`side_car_path` longtext,\n`side_car_hash` varchar(191) DEFAULT NULL,\n`blurhash` longtext,\nPRIMARY KEY (`id`),\nUNIQUE KEY `path_hash` (`path_hash`),\nUNIQUE KEY `side_car_hash` (`side_car_hash`),\nKEY `idx_media_video_metadata_id` (`video_metadata_id`),\nKEY `idx_media_album_id` (`album_id`),\nKEY `idx_media_exif_id` (`exif_id`),\nKEY `idx_media_type` (`type`),\nCONSTRAINT `fk_media_album` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_media_exif` FOREIGN KEY (`exif_id`) REFERENCES `media_exif` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_media_video_metadata` FOREIGN KEY (`video_metadata_id`) REFERENCES `video_metadata` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=880 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_media_urls_insert := "CREATE TABLE if not exists `media_urls` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`media_id` bigint(20) NOT NULL,\n`media_name` varchar(191) NOT NULL,\n`width` bigint(20) NOT NULL,\n`height` bigint(20) NOT NULL,\n`purpose` varchar(191) NOT NULL,\n`content_type` longtext NOT NULL,\n`file_size` bigint(20) NOT NULL,\nPRIMARY KEY (`id`),\nKEY `idx_media_urls_purpose` (`purpose`),\nKEY `idx_media_urls_media_id` (`media_id`),\nKEY `idx_media_urls_media_name` (`media_name`),\nCONSTRAINT `fk_media_media_url` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=427 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_user_media_data_insert := "CREATE TABLE if not exists `user_media_data` (\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`user_id` bigint(20) NOT NULL,\n`media_id` bigint(20) NOT NULL,\n`favorite` tinyint(1) NOT NULL DEFAULT '0',\nPRIMARY KEY (`user_id`,`media_id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci "
	sql_share_tokens_insert := "CREATE TABLE if not exists `share_tokens` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`value` longtext NOT NULL,\n`owner_id` bigint(20) NOT NULL,\n`expire` datetime(3) DEFAULT NULL,\n`password` longtext,\n`album_id` bigint(20) DEFAULT NULL,\n`media_id` bigint(20) DEFAULT NULL,\nPRIMARY KEY (`id`),\nKEY `idx_share_tokens_media_id` (`media_id`),\nKEY `idx_share_tokens_owner_id` (`owner_id`),\nKEY `idx_share_tokens_expire` (`expire`),\nKEY `idx_share_tokens_album_id` (`album_id`),\nCONSTRAINT `fk_share_tokens_album` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_share_tokens_media` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_share_tokens_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_face_groups_insert := "CREATE TABLE if not exists `face_groups` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`label` longtext,\nPRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	sql_image_faces_insert := "CREATE TABLE if not exists `image_faces` (\n`id` bigint(20) NOT NULL AUTO_INCREMENT,\n`created_at` datetime(3) DEFAULT NULL,\n`updated_at` datetime(3) DEFAULT NULL,\n`face_group_id` bigint(20) NOT NULL,\n`media_id` bigint(20) NOT NULL,\n`descriptor` blob NOT NULL,\n`rectangle` varchar(64) NOT NULL,\nPRIMARY KEY (`id`),\nKEY `idx_image_faces_face_group_id` (`face_group_id`),\nKEY `idx_image_faces_media_id` (`media_id`),\nCONSTRAINT `fk_face_groups_image_faces` FOREIGN KEY (`face_group_id`) REFERENCES `face_groups` (`id`) ON DELETE CASCADE,\nCONSTRAINT `fk_media_faces` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci"
	dataApi, _ := NewDataApiClient()
	//dataApi.ExecuteSQl(sql_site_info_insert)
	//dataApi.ExecuteSQl(sql_users_insert)
	//dataApi.ExecuteSQl(sql_access_tokens_insert)
	//dataApi.ExecuteSQl(sql_user_preferences_insert)
	//dataApi.ExecuteSQl(sql_albums_insert)
	//dataApi.ExecuteSQl(sql_user_albums_insert)
	//dataApi.ExecuteSQl(sql_video_metadata_insert)
	//dataApi.ExecuteSQl(sql_media_exif_insert)
	//dataApi.ExecuteSQl(sql_media_insert)
	//dataApi.ExecuteSQl(sql_media_urls_insert)
	//dataApi.ExecuteSQl(sql_user_media_data_insert)
	//dataApi.ExecuteSQl(sql_share_tokens_insert)
	//dataApi.ExecuteSQl(sql_face_groups_insert)
	//dataApi.ExecuteSQl(sql_image_faces_insert)
	var sql [14]string
	sql[0] = sql_site_info_insert
	sql[1] = sql_users_insert
	sql[2] = sql_access_tokens_insert
	sql[3] = sql_user_preferences_insert
	sql[4] = sql_albums_insert
	sql[5] = sql_user_albums_insert
	sql[6] = sql_video_metadata_insert
	sql[7] = sql_media_exif_insert
	sql[8] = sql_media_insert
	sql[9] = sql_media_urls_insert
	sql[10] = sql_user_media_data_insert
	sql[11] = sql_share_tokens_insert
	sql[12] = sql_face_groups_insert
	sql[13] = sql_image_faces_insert
	for i := 0; i < 14; i++ {
		dataApi.ExecuteSQl(sql[i])
	}
}
func Stresstest() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			dataApi, _ := NewDataApiClient()
			sql_serverless_test := "select benchmark(37000000 ,crc32('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef'))"
			res, _ := dataApi.Query(sql_serverless_test)
			println(res)
		}()
	}
	//for i := 0; i < 2; i++ {
	//	go func() {
	//		defer wg.Done()
	//		sql_serverless_test := "select benchmark(30000000 ,crc32('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef'))"
	//		dataApi.Query(sql_serverless_test)
	//	}()
	//}
	wg.Wait()
}
