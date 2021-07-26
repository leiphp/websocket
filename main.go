package main

import (
	"chatroom/configs"
	"chatroom/initialize"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zserge/lorca"
	"log"
	"net/http"
)

func main() {
	//获得配置对象
	Yaml := configs.InitConfig()
	initialize.Init(Yaml)

	router := mux.NewRouter()
	//监听channel数据
	go h.run()
	//开启chrome
	go StartChrome()
	// 配置 websocket route
	router.HandleFunc("/ws", myws)
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}

func StartChrome(){
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("E:\\github\\websocket\\local.html", "", 480, 320)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	// Wait until UI window is closed
	<-ui.Done()
}