package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/sayopaul/sendchamp-go-test/config"
)

type SocketService struct {
	configEnv config.Config
}

func NewSocketService(configEnv config.Config) SocketService {
	return SocketService{
		configEnv: configEnv,
	}
}

func (ss SocketService) Send(message map[string]interface{}) {
	var addr = flag.String("addr", fmt.Sprintf("%s:%s", ss.configEnv.SocketConnectionHost, ss.configEnv.SocketConnectionPort), "http service address")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Panic("dial:", err.Error())
	}
	defer c.Close()

	toJson, _ := json.Marshal(message)
	err = c.WriteMessage(websocket.TextMessage, []byte(toJson))
	if err != nil {
		log.Panic("error writing to the socket:", err.Error())
	}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read message error:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

}
