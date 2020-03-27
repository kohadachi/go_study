package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const URL = "http://weather.livedoor.com/forecast/webservice/json/v1?city=400040"

var DbConnection *sql.DB

type API struct {
	PinpointLocations PinpointLocations `json:pinpointLocations`
}

type PinpointLocation struct {
	Link string `json:link`
	Name string `json:name`
}

type PinpointLocations []*PinpointLocation

func Fetch() ([]byte, error) {
	// APIリクエストする
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(request)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	return responseBody, nil
}

func ImportJson(jsonByte []byte) (*API, error) {
	var data API
	if err := json.Unmarshal(jsonByte, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type Weather struct {
	gorm.Model
	Link string
	Name string
}

func main() {
	// ログ出力 (start)
	logger, _ := zap.NewDevelopment()
	logger.Info("start", zap.Time("date", time.Now()))

	// sqlite3にコネクト
	db, err := gorm.Open("sqlite3", "weather.db")
	if err != nil {
		logger.Error("start", zap.String("err", err.Error()), zap.Time("now", time.Now()))
	}
	defer db.Close()

	// マイグレーション
	db.AutoMigrate(&Weather{})
	// テーブル作成
	db.Table("weathers").CreateTable(&Weather{})

	jsonDataByte, _ := Fetch()
	convertData, _ := ImportJson(jsonDataByte)
	for _, pinpointLocation := range convertData.PinpointLocations {
		w := &Weather{Link: pinpointLocation.Link, Name: pinpointLocation.Name}
		db.Save(w)
	}

	// データ取得
	weather := Weather{}
	db.Table("weathers").Find(&weather, "id = ?", 2)
	fmt.Println(weather.Link)

	// ログ出力(end)
	logger.Info("end", zap.Time("date", time.Now()))
}
