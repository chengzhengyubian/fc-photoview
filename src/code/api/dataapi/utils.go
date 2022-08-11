package dataapi

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"rds-data-20220330/client"
	"reflect"
	"strings"
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
