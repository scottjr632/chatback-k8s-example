package main

import (
	"log"
	"net/http"

	"github.com/scottjr632/chatback/broker/sse"
)

const (
	port = ":50051"
)

func main() {
	sse := sse.New()
	http.HandleFunc("/subscribe", sse.Handler)
	http.HandleFunc("/message", sse.HandlePost)

	log.Printf("Listening on %v\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
