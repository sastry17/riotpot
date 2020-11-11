package mqttclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func connect(broker string, port int ) mqtt.Client {
	opts := createClientOptions(broker, port)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(broker string, port int) *mqtt.ClientOptions {

	pool := x509.NewCertPool()
	cert, err := tls.LoadX509KeyPair("./external/mqttclient/client.crt", "./external/mqttclient/client.key")
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{cert},
		// One-way authentication, that is, the client does not verify the server certificate.
		InsecureSkipVerify: true,
	}



	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker,port))
	opts.SetClientID("RIoTPot")
	opts.SetTLSConfig(tlsConfig)

	//opts.SetUsername(uri.User.Username())
	//password, _ := uri.User.Password()
	//opts.SetPassword(password)
	//opts.SetClientID(clientId)
	return opts
}



func Publisher(msg string) {
	var broker = "51.75.71.122"
	var port = 8883


	topic := "test"



	client := connect(broker, port)

	client.Publish(topic, 0, false, msg)

	}

