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
	"log"
	"os"

	"github.com/easygithdev/mqtt/client"
	"github.com/easygithdev/mqtt/client/conn"
)

func main() {

	// Show line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Disable log
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)

	clientId := "test-golang-mqtt"

	// connHost := "test.mosquitto.org"
	connHost := "mqtt.eclipseprojects.io"
	connPort := "1883"

	topic1 := "hello/mqtt1"
	topic2 := "hello/mqtt2"

	qos := client.QOS_0

	///////////////////////////////////////////////////////////
	// Init
	///////////////////////////////////////////////////////////
	mc := client.New(
		// client Id
		clientId,
		// connection infos
		client.WithConnInfos(conn.New(connHost, conn.WithPort(connPort))),
	)

	mc.OnSubscribe = func(mc client.MqttClient, userData interface{}, mid uint16) {
		fmt.Println("mid:", mid)
	}

	mc.OnMessage = func(mc client.MqttClient, userData interface{}, message string) {
		fmt.Println("msg: " + message)
	}

	// Connection

	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	///////////////////////////////////////////////////////////
	// Subscribe
	///////////////////////////////////////////////////////////

	resp1, err := mc.Subscribe(topic1, byte(qos))
	if err != nil || !resp1 {
		log.Printf("Subscribe Error: %s\n", err)
	}

	resp2, err := mc.Subscribe(topic2, byte(qos))
	if err != nil || !resp2 {
		log.Printf("Subscribe Error: %s\n", err)
	}

	if resp1 && resp2 {
		mc.LoopForever()
	}

}
