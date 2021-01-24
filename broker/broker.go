package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/scottjr632/chatback/broker/config"
	"github.com/scottjr632/chatback/broker/discovery"
	"github.com/scottjr632/chatback/broker/sse"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	port = ":50051"
)

func main() {
	godotenv.Load()
	configFile := os.Getenv("BROKER_CONFIG_PATH")
	if configFile == "" {
		configFile = "config.yml"
	}
	config, err := config.Load(configFile)
	if err != nil {
		panic(err)
	}

	disc, err := discovery.New(config)
	if err != nil {
		panic(err)
	}

	sse := sse.New(config, disc)
	http.HandleFunc("/subscribe", sse.Handler)
	http.HandleFunc("/message", sse.HandlePost)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on %v\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
