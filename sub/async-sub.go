package main

// Import Go and NATS packages
import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
)

func addGeolocation(client *redis.Client, data map[string]interface{}, msg *nats.Msg, id, countryCode, log, lat string) {
	if _, ok := data[id].(string); ok {
		err := client.GeoAdd(data[countryCode].(string),
			&redis.GeoLocation{Longitude: data[log].(float64), Latitude: data[lat].(float64), Name: string(msg.Data)}).Err()
		if err != nil {
			panic(err)
		}
	} else {
		err := client.GeoAdd(data[countryCode].(string),
			&redis.GeoLocation{Longitude: data[log].(float64), Latitude: data[lat].(float64), Name: string(msg.Data)}).Err()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Create server connection
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	log.Println("Connected to " + nats.DefaultURL)

	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "", // no password set
		// DB:       0,  // use default DB
	})

	// Subscribe to subject
	log.Printf("Subscribing to subject 'US-Bike-Sharing-Channel'\n")
	natsConnection.Subscribe("US-Bike-Sharing-Channel", func(msg *nats.Msg) {
		var dataRaw interface{}
		// Handle the message
		log.Printf("Received message '%s\n", string(msg.Data)+"'")

		json.Unmarshal(msg.Data, &dataRaw)

		data := dataRaw.(map[string]interface{})

		if data["station_id"] != nil {
			addGeolocation(client, data, msg, "station_id", "country_code", "lon", "lat")
		} else {
			addGeolocation(client, data, msg, "id", "country_code", "longitude", "latitude")
		}
	})

	// Keep the connection alive
	runtime.Goexit()
}
