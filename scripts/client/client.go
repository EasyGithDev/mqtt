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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/easygithdev/mqtt/client"
	"github.com/easygithdev/mqtt/client/conn"
)

func init() {

	// Show line numbers
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

var onConnect = func(mc client.MqttClient, userData interface{}, rc net.Conn) {
	fmt.Println("Connecting to server " + rc.RemoteAddr().String())
}

var onDisconnect = func(mc client.MqttClient, userData interface{}, rc net.Conn) {
	fmt.Println("Disconnect from server" + rc.RemoteAddr().String())
}

var onPublish = func(mc client.MqttClient, userData interface{}, mid uint16) {
	fmt.Printf("Publish\n")
}

var onSubscribe = func(mc client.MqttClient, userData interface{}, mid uint16) {
	fmt.Printf("Subscribe\n")
}

var onMessage = func(mc client.MqttClient, userData interface{}, message string) {
	fmt.Println("msg: " + message)
}

func idGenerator(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	pub := flag.Bool("pub", false, "publish")
	sub := flag.Bool("sub", false, "subscribe")
	host := flag.String("h", "", "the hostname")
	port := flag.String("p", "1883", "the port")
	id := flag.String("i", idGenerator(10), "the client Id")
	topic := flag.String("t", "", "the topic name")
	msg := flag.String("m", "", "the message to send")
	qos := flag.Uint("qos", 0, "quality of service")
	user := flag.String("user", "", "the username")
	pwd := flag.String("pwd", "", "the password")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	if !*pub && !*sub {
		fmt.Fprintf(os.Stderr, "missing required -pub/-sub flag\n")
		os.Exit(2)
	}

	if strings.TrimSpace(*host) == "" {
		fmt.Fprintf(os.Stderr, "missing required -h argument/flag\n")
		os.Exit(2)
	}

	if strings.TrimSpace(*topic) == "" {
		fmt.Fprintf(os.Stderr, "missing required -t argument/flag\n")
		os.Exit(2)
	}

	if *pub && strings.TrimSpace(*msg) == "" {
		fmt.Fprintf(os.Stderr, "missing required -m flag\n")
		os.Exit(2)
	}

	if !(*qos == 0 || *qos == 1 || *qos == 2) {
		fmt.Fprintf(os.Stderr, "qos must be 0,1 or 2\n")
		os.Exit(2)
	}

	if !*verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	var credentialsOption client.ClientOption = nil
	if *user != "" && *pwd != "" {
		credentialsOption = client.WithCredentials(*user, *pwd)
	}

	///////////////////////////////////////////////////////////
	// Init
	///////////////////////////////////////////////////////////
	mc := client.New(
		// client Id
		*id,
		// credentials
		credentialsOption,
		// connection infos
		client.WithConnInfos(conn.New(*host, conn.WithPort(*port))),
	)

	// Callbacks

	mc.OnConnect = onConnect
	mc.OnDisconnect = onDisconnect
	mc.OnPublish = onPublish
	mc.OnSubscribe = onSubscribe
	mc.OnMessage = onMessage

	// Connection

	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	///////////////////////////////////////////////////////////
	// Run
	///////////////////////////////////////////////////////////

	if *pub {
		resp, err := mc.Publish(*topic, *msg, byte(*qos))

		if err != nil || !resp {
			log.Printf("Publish Error: %s\n", err)
		}

	} else if *sub {

		resp, err := mc.Subscribe(*topic, byte(*qos))
		if err != nil {
			log.Printf("Subscribe Error: %s\n", err)
		}

		if resp {
			mc.LoopForever()
		}
	}

}
