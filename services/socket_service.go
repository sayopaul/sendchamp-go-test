package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

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

// Application constants, defining host, port, and protocol.
func (ss SocketService) Send(message map[string]interface{}) {
	// Start the client and connect to the server.
	log.Println("Connecting to a ", ss.configEnv.SocketConnectionType, "server with hostname", ss.configEnv.SocketConnectionHost+":"+ss.configEnv.SocketConnectionPort)
	//dial the server
	conn, err := net.Dial(ss.configEnv.SocketConnectionType, ss.configEnv.SocketConnectionHost+":"+ss.configEnv.SocketConnectionPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	//marshal to json
	toJson, _ := json.Marshal(message)
	//write message to connection
	conn.Write([]byte(toJson))
	// run loop forever, until the econnection is exited
	for {
		// Listen for response from server
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Print server response.
		log.Print("The server responded with: " + message)
	}
}
