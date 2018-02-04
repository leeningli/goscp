package main

import (
	"leeconfig"
	"fmt"
	"strings"
	"os"
	"net/http"
	"io"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
const (
	MAX_CNT_TOPIC = 100
)
var DB_TOPICS = [MAX_CNT_TOPIC]string

func readConfig() {
	fmt.Println("start read config field:main...")
	TOPIC := leeconfig.GetConfig("main")
	topics := TOPIC["topics"]
	topic_list := strings.Split(topics, ",")
	for k, db_topic := range topic_list {
		DB_TOPICS[k] = db_topic
	}
}

func init() {
	readConfig()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func pre_exporter() {
	for _, value := range DB_TOPICS {
		db_index := strings.Split(value, ":")
		if strings.ToLower(db_index) == "mysql" {
			mysql_exporter(value)
		} else if strings.ToLower(db_index) == "oracle" {
			oracle_exporter(value)
		} else {
			fmt.Println("config:", value, " is error.")
			os.Exit(1)
		}
	}
}

func mysql_exporter(appname string) {
	TOPIC := leeconfig.GetConfig(appname)
	ip := TOPIC["ip"]
	port := TOPIC["port"]
	username := TOPIC["username"]
	pwd := TOPIC["pwd"]
	dbname := TOPIC["dbname"]
	cmd := TOPIC["cmd"]
	sid := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, pwd, ip, port, dbname)
	fmt.Println("sid==", sid)
	db ,err := sql.Open("mysql", sid)
	checkError(err)
	rows, err := db.Query(cmd)
	checkError(err)
}

func oracle_exporter(appname string) {
	TOPIC := leeconfig.GetConfig(appname)
	ip := TOPIC["ip"]
	port := TOPIC["port"]
	username := TOPIC["username"]
	pwd := TOPIC["pwd"]
}

func ExporterHandler(w http.ResponseWriter, r *http.Request) {
	var result := `flag {flag="b",flag2="a"} 123`
	io.WriteString(w, result)
}
