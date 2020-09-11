package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	sfu "github.com/pion/ion-sfu/pkg
	webrtc "github.com/pion/webrtc/v3"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	server *socketio.Server
	conf   = sfu.Config{}
	file   string
)
const (
	portRangeLimit = 100
)

func setConfigs() bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file %s read failed. %v\n", file, err)
		return false
	}
	err = viper.GetViper().Unmarshal(&conf)
	if err != nil {
		fmt.Printf("sfu config file %s loaded failed. %v\n", file, err)
		return false
	}

	if len(conf.WebRTC.ICEPortRange) > 2 {
		fmt.Printf("config file %s loaded failed. range port must be [min,max]\n", file)
		return false
	}

	if len(conf.WebRTC.ICEPortRange) != 0 && conf.WebRTC.ICEPortRange[1]-conf.WebRTC.ICEPortRange[0] < portRangeLimit {
		fmt.Printf("config file %s loaded failed. range port must be [min, max] and max - min >= %d\n", file, portRangeLimit)
		return false
	}

	fmt.Printf("config %s load ok!\n", file)
	return true
}

type Join struct {
	Sid   string                    `json:"sid"`
	Offer webrtc.SessionDescription `json:"offer"`
}

func startSocket() {
	_server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server = _server
}
func setSocketHandlers() {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {

	})
	server.OnEvent("/", "offer", func(s socketio.Conn, msg string) {

	})
	server.OnEvent("/", "answer", func(s socketio.Conn, msg string) {

	})
	server.OnEvent("/", "trickle", func(s socketio.Conn, msg string) {

	})
}

func main() {

}
