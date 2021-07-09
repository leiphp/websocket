package nats

import (
	"chatroom/infra/database"
	"chatroom/infra/vars"
	"chatroom/initialize"
	"encoding/json"
	"fmt"
	"log"
)

//向ui推送外卖通知
func NotifyTakeout(mgs *database.Data) {
	subject := fmt.Sprintf(vars.TakeoutNotifySubject, "123456789")
	data, err := json.Marshal(mgs)
	if err != nil {
		log.Println("Marshal err", err.Error())
	}
	log.Println("ui takeout message:", string(data))
	err = initialize.NatsClient.Publish(subject, data)
	if err != nil {
		log.Println("Nats Publish err", err.Error())
	}
}
