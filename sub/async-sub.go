package main

// Import Go and NATS packages
import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/go-redis/redis"
	"github.com/nats-io/go-nats"
)

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
			if _, ok := data["station_id"].(string); ok {
				err := client.GeoAdd(data["country_code"].(string),
					&redis.GeoLocation{Longitude: data["lon"].(float64), Latitude: data["lat"].(float64), Name: string(msg.Data)}).Err()
				if err != nil {
					panic(err)
				}
			} else {
				err := client.GeoAdd(data["country_code"].(string),
					&redis.GeoLocation{Longitude: data["lon"].(float64), Latitude: data["lat"].(float64), Name: string(msg.Data)}).Err()
				if err != nil {
					panic(err)
				}
			}
		} else {
			if _, ok := data["id"].(string); ok {
				err := client.GeoAdd(data["country_code"].(string),
					&redis.GeoLocation{Longitude: data["longitude"].(float64), Latitude: data["latitude"].(float64), Name: string(msg.Data)}).Err()
				if err != nil {
					panic(err)
				}
			} else {
				err := client.GeoAdd(data["country_code"].(string),
					&redis.GeoLocation{Longitude: data["longitude"].(float64), Latitude: data["latitude"].(float64), Name: string(msg.Data)}).Err()
				if err != nil {
					panic(err)
				}
			}
		}
	})

	// Keep the connection alive
	runtime.Goexit()
}
