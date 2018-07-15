package main

import (
	"fmt"
	"log"
	"os"

	nats "github.com/nats-io/go-nats"
	"github.com/satori-com/satori-rtm-sdk-go/rtm"
	"github.com/satori-com/satori-rtm-sdk-go/rtm/pdu"
	"github.com/satori-com/satori-rtm-sdk-go/rtm/subscription"
)

// Import packages

const (
	ENDPOINT = "wss://open-data.api.satori.com"
	// APP_KEY  = "b78adebf97eBe26FDcD4CD1B82c56a2f"
	CHANNEL = "US-Bike-Sharing-Channel"
	SUBJECT = "US-Bike-Sharing-Channel"
)

func main() {
	//GET APP KEY
	appKey := os.Args[1]

	// Connect to server; defer close
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	defer natsConnection.Close()
	log.Println("Connected to " + nats.DefaultURL)

	client, err := rtm.New(ENDPOINT, appKey, rtm.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client.OnConnected(func() {
		fmt.Println("Connected to Satori RTM!")
	})

	data_c := make(chan string)
	listener := subscription.Listener{
		OnData: func(data pdu.SubscriptionData) {
			for _, message := range data.Messages {
				data_c <- string(message)
			}
		},
	}
	client.Subscribe(CHANNEL, subscription.SIMPLE, pdu.SubscribeBodyOpts{}, listener)

	client.Start()

	for message := range data_c {
		// fmt.Println("Got message:", message)
		natsConnection.Publish(SUBJECT, []byte(message))
	}

	// log.Println("Published message on subject " + subject)
}
