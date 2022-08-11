package dataapi

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/photoview/photoview/api/utils"
)
import "rds-data-20220330/client"

//初始化客户端配置

type DataApi struct {
	client   *client.Client
	resqust  *client.ExecuteStatementRequest
	respnose client.ExecuteStatementResponse
	res      [][]*client.ExecuteStatementResponseBodyDataRecords
}

func NewDataApiClient() (*DataApi, error) {
	accessKeyId := utils.AccessKeyId.GetValue()
	accessKeySecret := utils.AccessKeySecret.GetValue()
	endpoint := utils.Endpoint.GetValue()
	database := utils.Database.GetValue()
	secretArn := utils.SecretArn.GetValue()
	resource := utils.ResourceArn.GetValue()
	config := openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		Endpoint:        &endpoint,
	}
	dataClient, _ := client.NewClient(&config)
	dataRequst := client.ExecuteStatementRequest{
		Database:    &database,
		ResourceArn: &resource,
		SecretArn:   &secretArn,
	}

	return &DataApi{
		client:  dataClient,
		resqust: &dataRequst,
	}, nil
}

//func (dataApi *DataApi) GenSql(sqlTemplate string, args string[]) (string, error) {
//	num := len(sqlTemplate)
//	for i := 0; i < num; i++ {
//
//	}
//}

func NewDataApiClientJosn() (*DataApi, error) {
	var formatRecordsAs string
	formatRecordsAs = "JSON"
	accessKeyId := utils.AccessKeyId.GetValue()
	accessKeySecret := utils.AccessKeySecret.GetValue()
	endpoint := utils.Endpoint.GetValue()
	database := utils.Database.GetValue()
	secretArn := utils.SecretArn.GetValue()
	resource := utils.ResourceArn.GetValue()
	config := openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		Endpoint:        &endpoint,
	}
	dataClient, _ := client.NewClient(&config)
	dataRequst := client.ExecuteStatementRequest{
		Database:        &database,
		ResourceArn:     &resource,
		SecretArn:       &secretArn,
		FormatRecordsAs: &formatRecordsAs,
	}
	return &DataApi{
		client:  dataClient,
		resqust: &dataRequst,
	}, nil
}

func (dataApi *DataApi) ExecuteSQl(sql string) (*client.ExecuteStatementResponse, error) {
	dataApi.resqust.Sql = &sql
	req := &dataApi.respnose
	req, err := dataApi.client.ExecuteStatement(dataApi.resqust)
	if err != nil {
		fmt.Println(err)
	}
	if req == nil {
		return nil, nil
	}
	return req, nil
}

func (dataApi *DataApi) Query(sql string) ([][]*client.ExecuteStatementResponseBodyDataRecords, error) {
	dataApi.resqust.Sql = &sql
	req := &dataApi.respnose
	req, err := dataApi.client.ExecuteStatement(dataApi.resqust)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return req.Body.Data.Records, nil
}
