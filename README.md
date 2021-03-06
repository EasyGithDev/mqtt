# Mqtt client
I'll try to implement the draft as mentionned here :

https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/


## Use as a binary

Open a terminal and write :
```bash
    ./make
```

Or you can edit the make file to perfome your own buid command.

You will find the binary in the bin folder.

You can send a publish like above :
```bash
    ./bin/client -pub -h=test.mosquitto.org -t=hello/mqtt "-m=hello world"
```

You can subscribe to a topic and read the messages like above :
```bash
    ./bin/client -sub -h=test.mosquitto.org -t=hello/mqtt
```

## Use as a library

#### Connection

Simple connect with :
- clientId
- hostname
- port

```go
        clientId := "test-golang-mqtt"
        connHost := "test.mosquitto.org"
        connPort := "1883"

        mc := client.New(
            // client Id
            clientId,
            // connection infos
            client.WithConnInfos(conn.New(connHost, conn.WithPort(connPort))),
        )

        _, connErr := mc.Connect()
        if connErr != nil {
            log.Print("Error connecting:", connErr.Error())
        }
        defer mc.Close()
```

Connect with credentials :
- Username
- Password

```go

        ...

        username := "rw"
        password := "readwrite"
        mc := client.New(
            // client Id
            clientId,
            // Credentials
            client.WithCredentials(username, password),
            // connection infos
            client.WithConnInfos(conn.New(connHost, conn.WithPort(connPort))),
        )

        ...
	
```

#### Publish

Publish a message :
- topic : the topic that the message should be published on.
- message : the actual message to send.
- qos : the quality of service level to use.
- retain : if set to True, the message will be set as the “last known good”/retained message for the topic.

```go
        topic := "hello/mqtt"
        qos := client.QOS_0
        msg := "The temperature is 5 degrees"

        _, pubErr := mc.Publish(topic, msg, byte(qos))

        if pubErr != nil {
            log.Print("Error publishing:", pubErr.Error())
        }
```

Publish many messages :

```go

        go mc.LoopStart()

        for {
            temperature := rand.Intn(60)
            msg := "The temperature is " + fmt.Sprintf("%d", temperature)
            _, pubErr := mc.Publish(topic, msg, byte(qos))

            if pubErr != nil {
                log.Print("Error publishing:", pubErr.Error())
                break
            }
            time.Sleep(5 * time.Second)
        }

```

#### Subscribe


Subscribe with :
- topic : a string specifying the subscription topic to subscribe to.
- qos : the desired quality of service level for the subscription.

```go
        topic := "hello/mqtt"
        qos := client.QOS_0

        _, errSub := mc.Subscribe(topic, client.QOS_0)
        if errSub != nil {
            log.Printf("Subscribe Error: %s\n", errSub)
        }
```

Get the messages :

```go
        _, errSub := mc.Subscribe(topic, client.QOS_0)
        if errSub != nil {
            log.Printf("Subscribe Error: %s\n", errSub)
        } else {
            mc.LoopForever()
        }
```

#### Callback function

```go
        var onConnect = func(mc client.MqttClient, userData interface{}, rc net.Conn) {
        fmt.Println("Connecting to server " + rc.RemoteAddr().String())
        }

        var onDisconnect = func(mc client.MqttClient, userData interface{}, rc net.Conn) {
        fmt.Println("Disconnect from server" + rc.RemoteAddr().String())
        }

        var onPublish = func(mc client.MqttClient, userData interface{}, mid int) {
        fmt.Printf("Publish\n")
        }

        var onSubscribe = func(mc client.MqttClient, userData interface{}, mid int) {
        fmt.Printf("Subscribe\n")
        }

        var onMessage = func(mc client.MqttClient, userData interface{}, message string) {
        fmt.Println("msg: " + message)
        }

        ...

        mc.OnConnect = onConnect
        mc.OnDisconnect = onDisconnect
        mc.OnPublish = onPublish
        mc.OnSubscribe = onSubscribe
        mc.OnMessage = onMessage
```


