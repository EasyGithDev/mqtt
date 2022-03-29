package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/easygithdev/mqtt/client"
)

const (
	// connHost = "localhost"
	// connPort = "8080"
	// connType = "tcp"

	// connHost = "test.mosquitto.org"
	connHost = "mqtt.eclipseprojects.io"

	connPort = "1883"
	// connPort = "1884"
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

	///////////////////////////////////////////////////////////
	// Connect
	///////////////////////////////////////////////////////////
	mc := client.NewMqttClient("go-lang-mqtt")

	_, err := mc.Connect(connHost, connPort)
	if err != nil {
		log.Print("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer mc.Disconnect()

	///////////////////////////////////////////////////////////
	// Ping
	///////////////////////////////////////////////////////////

	resp, err := mc.Ping()

	if err != nil {
		log.Print("Error ping:", err.Error())
	}

	if resp {
		log.Print("Ping is done with success")

	}

	// log.Print("Connecting to " + connType + " server " + connHost + ":" + connPort)

	///////////////////////////////////////////////////////////
	// Publish
	///////////////////////////////////////////////////////////

	// mc.SetOptions(&client.MqttConnectOptions{Login: "aa", Password: "bb"})

	// _, err = mc.Publish("/hello/world", "this is my hello world")

	// if err != nil {
	// 	log.Print("Error publishing:", err.Error())
	// }

	///////////////////////////////////////////////////////////
	// Subscribe
	///////////////////////////////////////////////////////////

	// resp, err := mc.Subscribe("/tartine/de/confiture")
	// if err != nil {
	// 	log.Printf("Subscribe Error: %s\n", err)
	// }

	// if resp {
	// 	log.Printf("Subcribe established \n")
	// 	mc.ReadLoop()
	// }

}
