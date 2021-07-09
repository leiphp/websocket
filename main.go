package main

import (
	"chatroom/configs"
	"chatroom/initialize"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	//获得配置对象
	Yaml := configs.InitConfig()
	initialize.Init(Yaml)

	router := mux.NewRouter()
	//监听channel数据
	go h.run()
	// 配置 websocket route
	router.HandleFunc("/ws", myws)
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}