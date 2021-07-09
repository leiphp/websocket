package main

import (
	."chatroom/infra/database"
	"chatroom/pkg/nats"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//为什么用指针，因为在连接对象时这些值可以都已经赋值了
type connection struct {
	ws   *websocket.Conn
	sc   chan []byte
	data *Data
}

var wu = &websocket.Upgrader{
	ReadBufferSize: 512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func myws(w http.ResponseWriter, r *http.Request) {
	// 将初始GET请求升级到websocket
	ws, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{sc: make(chan []byte, 256), ws: ws, data: &Data{}}
	h.r <- c //向register连接器注册请求
	go c.writer()
	c.reader()
	defer func() {
		c.data.Type = "logout"
		user_list = del(user_list, c.data.User)
		c.data.UserList = user_list
		c.data.Content = c.data.User
		data_b, _ := json.Marshal(c.data)
		h.b <- data_b
		h.r <- c
	}()
}

//哪个client需要写入哪个client就调用
func (c *connection) writer() {
	for message := range c.sc {
		m := make(map[string]interface{})
		json.Unmarshal(message, &m)
		log.Println("writer message:",m)
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

var user_list = []string{}

//读取客户端连接ws发送的数据
func (c *connection) reader() {
	for {
		m := make(map[string]interface{})
		_, message, err := c.ws.ReadMessage()
		json.Unmarshal(message, &m)
		log.Println("reader message:",m)

		if err != nil {
			h.r <- c
			break
		}
		json.Unmarshal(message, &c.data)
		log.Println("reader data:",c.data)
		//读到消息通过nats推送到ui界面
		go nats.NotifyTakeout(c.data)
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content
			c.data.From = c.data.User
			user_list = append(user_list, c.data.User)
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "user":
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "logout":
			c.data.Type = "logout"
			user_list = del(user_list, c.data.User)
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
			h.r <- c
		default:
			fmt.Print("========default================")
		}
	}
}

//删除用户切片中数据
func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}
