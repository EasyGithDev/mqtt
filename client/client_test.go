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
package client

import (
	"log"
	"os"
	"testing"

	"github.com/easygithdev/mqtt/client/conn"
)

const (
	connHost = "mqtt.eclipseprojects.io"
	connPort = "1883"
)

func TestConnectDisconnect(t *testing.T) {

	mc := New("test-golang-mqtt",
		WithConnInfos(conn.New(connHost)),
	)
	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	response, err := mc.MqttConnect()

	if err != nil {
		t.Errorf("Mqtt connection fail %s", err)
	} else if !response {
		t.Errorf("Mqtt connection fail")
	} else {

		response, err := mc.MqttDisconnect()

		if err != nil {
			t.Errorf("Mqtt disconnect fail %s", err)
		} else if !response {
			t.Errorf("Mqtt disconnect fail")
		}
	}
}

func TestPing(t *testing.T) {

	mc := New("test-golang-mqtt",
		WithConnInfos(conn.New(connHost)),
	)
	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	response, err := mc.Ping()

	if err != nil {
		t.Errorf("Mqtt ping fail %s", err)
	}
	if !response {
		t.Errorf("Mqtt ping fail")
	}

}

func TestPublish(t *testing.T) {

	mc := New("test-golang-mqtt",
		WithConnInfos(conn.New(connHost)),
	)
	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	response, err := mc.Publish("/hello/world", "this is my hello world", 0)

	if err != nil {
		t.Errorf("Mqtt publish fail %s", err)
	}
	if !response {
		t.Errorf("Mqtt publish fail")
	}

}

func TestSubscribe(t *testing.T) {

	mc := New("test-golang-mqtt",
		WithConnInfos(conn.New(connHost)),
	)
	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	response, err := mc.Subscribe("/hello/world", 0)

	if err != nil {
		t.Errorf("Mqtt subscribe fail %s", err)
	}

	if !response {
		t.Errorf("Mqtt subscribe fail")
	}

}
