package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/easygithdev/mqtt/client"
)

const (
	// connHost = "localhost"
	// connPort = "8080"
	// connType = "tcp"

	connHost = "test.mosquitto.org"
	//connPort = "1883"
	connPort = "1884"
	connType = "tcp"

	DELIMITER byte = '\n'
	QUIT_SIGN      = "quit!"

	DEBUG = true
)

func init() {

	if !DEBUG {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

}

func main() {

	conn, err := net.Dial(connType, connHost+":"+connPort)

	if err != nil {
		log.Print("Error connecting:", err.Error())
		os.Exit(1)
	}
	log.Print("Connecting to " + connType + " server " + connHost + ":" + connPort)

	// Connect
	mc := client.NewMqttClient()

	// Adding connection to mc
	mc.SetConn(&conn)
	options := &client.MqttConnectOptions{Login: "rw", Password: "readwrite"}
	//options = nil
	resp, err := mc.Connect("golang.test-1", options)

	if err != nil {
		log.Printf("Connect Error: %s\n", err)
	}

	if resp {
		log.Printf("Connection established: \n")
	}

	// resp, err = mc.Disconnect()

	// if err != nil {
	// 	log.Printf("Connect Error: %s\n", err)
	// }

	// if resp {
	// 	log.Printf("Disconnected: \n")
	// }

}
