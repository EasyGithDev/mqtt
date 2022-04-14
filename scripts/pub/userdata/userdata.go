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
	"math/rand"
	"os"
	"time"

	"github.com/easygithdev/mqtt/client"
	"github.com/easygithdev/mqtt/client/conn"
)

func main() {

	clientId := "test-golang-mqtt"

	connHost := "test.mosquitto.org"
	// connHost := "mqtt.eclipseprojects.io"

	// connPort := "1883"
	connPort := "1884"

	topic := "hello/mqtt"

	qos := client.QOS_2

	var nbpub int = 0

	// Show line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Disable log
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	///////////////////////////////////////////////////////////
	// Init
	///////////////////////////////////////////////////////////
	mc := client.New(
		// client Id
		clientId,
		client.WithCredentials("rw", "readwrite"),
		// connection infos
		client.WithConnInfos(conn.New(connHost, conn.WithPort(connPort))),
		// user data
		client.WithUserData(&nbpub),
	)
	mc.OnPublish = func(mc client.MqttClient, userData interface{}, mid uint16) {
		(*userData.(*int))++
		fmt.Println("Publish performed with mid:", mid)
	}

	// Connection

	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	///////////////////////////////////////////////////////////
	// Publish Loop
	///////////////////////////////////////////////////////////

	go mc.LoopStart()

	for {
		temperature := rand.Intn(60)
		msg := "The temperature is " + fmt.Sprintf("%d", temperature)
		resp, err := mc.Publish(topic, msg, byte(qos))

		if err != nil || !resp {
			log.Print("Error publishing:", err)
			continue
		}
		time.Sleep(2 * time.Second)

		fmt.Println("Nb publish:", nbpub)
		// Stop after 10 publish
		if nbpub == 10 {
			break
		}
	}

	fmt.Println("Publish performed:", nbpub)

}
