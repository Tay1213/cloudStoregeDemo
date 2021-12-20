package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"time"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var conf *ini.File

func Setup() {
	var err error
	conf, err = ini.Load("conf/app.ini")
	if err != nil {
		fmt.Printf("出错了:%#v\n", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("server", ServerSetting)
	mapTo("redis", RedisSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

}

func mapTo(s string, v interface{}) {
	err := conf.Section(s).MapTo(v)
	if err != nil {
		fmt.Printf("出错了: %#v\n", err)
	}
}
