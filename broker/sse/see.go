package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/scottjr632/chatback/broker/config"
	"github.com/scottjr632/chatback/broker/discovery"
	"github.com/scottjr632/chatback/broker/metrics"
)

// SSE handles the sse connections
type SSE interface {
	Handler(w http.ResponseWriter, r *http.Request)
	HandlePost(w http.ResponseWriter, r *http.Request)
}

type sse struct {
	clients   sync.Map
	config    *config.Config
	discovery *discovery.Discovery

	peers      chan *message
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
func New(config *config.Config, discovery *discovery.Discovery) SSE {
	sse := &sse{
		clients:    sync.Map{},
		config:     config,
		discovery:  discovery,
		peers:      make(chan *message),
		broadcast:  make(chan *message),
		register:   make(chan *client),
		unRegister: make(chan *client),
	}
	go sse.startListener()
	if config.Discovery {
		go sse.startPeerBroadcaster()
	}
	return sse
}

// Handler handles a new SSE connection
func (s *sse) Handler(w http.ResponseWriter, r *http.Request) {
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

	s.register <- client
	defer func() { s.unRegister <- client }()

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
func (s *sse) HandlePost(w http.ResponseWriter, r *http.Request) {
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

	s.broadcast <- message
	if s.config.Discovery && r.URL.Query().Get("propagate") == "" {
		s.peers <- message
	}
	log.Printf("got new message: %v\n", message)
	w.WriteHeader(http.StatusOK)
}

func (s *sse) startListener() {
	log.Println("Starting listener")
	for {
		select {
		case client := <-s.register:
			log.Println("Got new client")
			metrics.Get().IncActiceClients()
			s.clients.Store(client, true)
		case client := <-s.unRegister:
			log.Println("Unregistering client")
			metrics.Get().DecActiveClients()
			s.clients.Delete(client)
		case msg := <-s.broadcast:
			log.Printf("Broadcasting message: %v\n", msg)
			metrics.Get().IncMessagesSent()
			s.clients.Range(func(k interface{}, value interface{}) bool {
				client := k.(*client)
				client.messageChan <- msg
				return true
			})
		}
	}
}

func (s *sse) startPeerBroadcaster() {
	client := http.Client{}
	myIPs, err := getMachineIPs()
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := <-s.peers
		ips, err := s.discovery.GetBrokerPodsIPs()
		if err != nil {
			log.Println(err)
			continue
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			continue
		}

		msgReader := bytes.NewBuffer(msgJSON)
		for _, ip := range ips {
			if myIPs[ip] {
				continue
			}

			// we need to add propagate=false so the broker does not resend a message
			url := fmt.Sprintf("http://%s:%s/message?propagate=false", ip, s.config.Port)
			resp, err := client.Post(url, "application/json", msgReader)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("status_code: " + resp.Status)
		}
	}
}

func getMachineIPs() (map[string]bool, error) {
	myIPs := map[string]bool{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return myIPs, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return myIPs, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			myIPs[ip.String()] = true
		}
	}
	return myIPs, err
}
