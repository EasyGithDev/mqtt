# Mqtt client
Implement MQTT client in Golang.

Work in progess ...

## use as a library

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
- topic
- qos
- message

```go
        topic := "hello/mqtt"
    	qos := client.QOS_0
        msg := "The temperature is 5 degrees"

	    _, pubErr := mc.Publish(topic, msg, byte(qos))

		if pubErr != nil {
			log.Print("Error publishing:", pubErr.Error())
		}
```

#### Subscribe


Subscribe with :
- topic
- qos

```go
        topic := "hello/mqtt"
        qos := client.QOS_0

        _, errSub := mc.Subscribe(topic, client.QOS_0)
        if errSub != nil {
            log.Printf("Subscribe Error: %s\n", errSub)
        }

```

## use as a binary


