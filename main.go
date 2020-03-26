package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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

func main() {
	var data API
	// ログ出力 (start)
	fmt.Println("start")

	// APIリクエストする
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.DefaultClient.Do(request)
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(responseBody, &data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	for _, post := range data.PinpointLocations {
		fmt.Println(post.Link)
		fmt.Println(post.Name)
	}

	// 取得したデータをDBに保存する
	DbConnection, _ := sql.Open("sqlite3", ".weather.sql")
	cmd := `CREATE TABLE IF NOT EXISTS weather(
		link STRING,
		name STRING)
		`
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	// ログ出力(end)
}
