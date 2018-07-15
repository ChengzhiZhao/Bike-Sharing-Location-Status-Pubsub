# Bike-Sharing-Location-Status-Pubsub
Bike Sharing Location Status Streaming in NATS

Quick Start

### 1. Install nats-io/gnatsd
https://github.com/nats-io/gnatsd#quickstart

### 2. Start gnatsd server 
./gnatsd

### 3. Start Publisher
go run pub/pub-simple.go APP_KEY

### 4. Start Subscriber (also need redis to persist data or you can remove the line to post to redis) 
go run async-sub.go
