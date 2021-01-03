package sse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/scottjr632/chatback/broker/metrics"
)

// SSE handles the sse connections
type SSE interface {
	Handler(w http.ResponseWriter, r *http.Request)
	HandlePost(w http.ResponseWriter, r *http.Request)
}

type sse struct {
	clients sync.Map

	broadcast  chan *message
	register   chan *client
	unRegister chan *client
}

type client struct {
	messageChan chan *message
}

type message struct {
	ID   []byte `json:"id"`
	Data []byte `json:"data"`
}

// New creates a new SSE
func New() SSE {
	sse := &sse{
		clients:    sync.Map{},
		broadcast:  make(chan *message),
		register:   make(chan *client),
		unRegister: make(chan *client),
	}
	go sse.startListener()
	return sse
}

// Handler handles a new SSE connection
func (sse *sse) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("got a new conn")
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Unable to create flusher", http.StatusInternalServerError)
		return
	}

	c, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "Unable to create close notifier", http.StatusInternalServerError)
		return
	}

	client := &client{
		messageChan: make(chan *message),
	}

	sse.register <- client
	defer func() { sse.unRegister <- client }()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-client.messageChan:
			fmt.Fprintf(w, "id: %s\n", msg.ID)
			fmt.Fprintf(w, "data: %s\n\n", msg.Data)
			f.Flush()
		case <-c.CloseNotify():
			return
		}
	}
}

// HandlePost handles a new message and sends the message to all attached
// clients
func (sse *sse) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	message := &message{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Unable to read content", http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, message); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	sse.broadcast <- message
	log.Printf("got new message: %v\n", message)
	w.WriteHeader(http.StatusOK)
}

func (sse *sse) startListener() {
	log.Println("Starting listener")
	for {
		select {
		case client := <-sse.register:
			log.Println("Got new client")
			metrics.Get().IncActiceClients()
			sse.clients.Store(client, true)
		case client := <-sse.unRegister:
			log.Println("Unregistering client")
			metrics.Get().DecActiveClients()
			sse.clients.Delete(client)
		case msg := <-sse.broadcast:
			log.Printf("Broadcasting message: %v\n", msg)
			metrics.Get().IncMessagesSent()
			sse.clients.Range(func(k interface{}, value interface{}) bool {
				client := k.(*client)
				client.messageChan <- msg
				return true
			})
		}
	}
}
