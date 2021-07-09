package initialize

import (
	"github.com/nats-io/go-nats"
	"github.com/spf13/viper"
)

var (
	NatsClient   *nats.Conn //nats客户端
	Config       *viper.Viper      //全局配置
)

//	提供系统初始化，全局变量
func Init(config *viper.Viper) {
	//nats初始化
	Config = config
	nc, err := nats.Connect(config.GetString("Nats.Url"))
	if err != nil {
		panic(err)
	}
	NatsClient = nc
}
