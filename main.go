// MIT License

// Copyright (c) 2022 Florent Brusciano

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/easygithdev/mqtt/client"
)

const (

	// connHost = "test.mosquitto.org"
	connHost = "mqtt.eclipseprojects.io"

	connPort = "1883"
	// connPort = "1884"

	DEBUG = false
)

func init() {

	if !DEBUG {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

}

var onTcpConnect = func(conn net.Conn) {
	fmt.Println("Connecting to server " + conn.RemoteAddr().String())
}

var onTcpDisconnect = func() {
	fmt.Println("Disconnect from server")
}

var onConnect = func() {
	fmt.Println("Connecting to MQTT server ...")
}

var onPublish = func() {
	fmt.Println("Publish ")
}

var onSubscribe = func() {
	fmt.Println("Subscribed ...")
}

var onMessage = func(message string) {
	fmt.Println("Message: " + message)
}

func main() {

	///////////////////////////////////////////////////////////
	// Connect
	///////////////////////////////////////////////////////////
	mc := client.NewMqttClient("go-lang-mqtt")
	mc.OnTcpConnect = onTcpConnect
	mc.OnTcpDisconnect = onTcpDisconnect
	mc.OnConnect = onConnect
	mc.OnPublish = onPublish
	mc.OnSubscribe = onSubscribe
	mc.OnMessage = onMessage

	_, connErr := mc.TcpConnect(connHost, connPort)
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.TcpDisconnect()

	///////////////////////////////////////////////////////////
	// Ping
	///////////////////////////////////////////////////////////

	// resp, err := mc.Ping()

	// if err != nil {
	// 	log.Print("Error ping:", err.Error())
	// }

	// if resp {
	// 	log.Print("Ping is done with success")
	// }

	// log.Print("Connecting to " + connType + " server " + connHost + ":" + connPort)

	///////////////////////////////////////////////////////////
	// Options for using username&password
	///////////////////////////////////////////////////////////

	mc.SetOptions(&client.MqttConnectOptions{Login: "aa", Password: "bb"})

	///////////////////////////////////////////////////////////
	// Publish
	///////////////////////////////////////////////////////////

	// pubResp, pubErr := mc.Publish("/hello/world", "this is my hello world")

	// if pubErr != nil {
	// 	log.Print("Error publishing:", pubErr.Error())
	// }

	// if pubResp {
	// 	log.Print("Publish is done with success")
	// }

	///////////////////////////////////////////////////////////
	// Subscribe
	///////////////////////////////////////////////////////////

	respSub, errSub := mc.Subscribe("/tartine/de/confiture")
	if errSub != nil {
		log.Printf("Subscribe Error: %s\n", errSub)
	}

	if respSub {
		log.Printf("Subcribe established \n")
		mc.ReadLoop()
	}

}
