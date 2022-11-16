package main

import (
	"fmt"
	"net/http"
	"xshortUrl/app"
	"xshortUrl/app/config"
	"xshortUrl/app/db"

	"github.com/spf13/viper"
)

func init() {
	config.Init()

	d := db.NewDB(
		viper.GetString("mysql.user"),
		viper.GetString("mysql.pwd"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.dbname"),
	)
	db.InitConn(d)
}

func main() {
	fmt.Println("短链服务启动中...")
	mux := http.NewServeMux()

	// 路由设置
	mux.HandleFunc("/", app.LogMiddleware(app.ReIndex))
	// 获取短链
	mux.HandleFunc("/getShortUrl", app.GetShortUrl)
	// 解析短链
	mux.HandleFunc("/parse", app.Parse)
	// 初始页面
	mux.HandleFunc("/index", app.Home)

	server := http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: mux,
	}

	fmt.Println("短链服务启动完成...")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
