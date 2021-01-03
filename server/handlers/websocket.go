package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/scottjr632/chatback/server/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type client struct{}

var clients = make(map[*websocket.Conn]client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan models.Message, 10)
var unRegister = make(chan *websocket.Conn)

func sendEventToClients(event models.Message) {
	for client := range clients {
		if err := client.WriteJSON(event); err != nil {
			log.Println("writeJson: ", err)
		}
	}
}

func (h *Handlers) handleRegisters() {
	sseChan, errChan := h.broker.HandleSubscribe()
	for {
		select {
		case newClient := <-register:
			log.Println("got new client")
			clients[newClient] = client{}
		case remClient := <-unRegister:
			log.Println("unregistered new client")
			delete(clients, remClient)
		case event := <-broadcast:
			sendEventToClients(event)
		case err := <-errChan:
			log.Printf("Error HandleSubscribe: %v\n", err)
		case event := <-sseChan:
			message := models.Message{}
			if err := json.Unmarshal(event.Data, &message); err != nil {
				log.Printf("Err ssechan: %v", err)
				continue
			}
			sendEventToClients(message)
		}
	}
}

func (*Handlers) UpgradeWebsocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handlers) WebsocketMessages(c *websocket.Conn) {
	register <- c
	defer func() {
		unRegister <- c
		c.Close()
	}()

	for {
		msg := models.Message{}
		if err := c.ReadJSON(&msg); err != nil {
			log.Println("read: ", err)
			break
		}
		log.Printf("recv: %v", msg)
		if err := msg.Insert(context.TODO(), h.db, boil.Infer()); err != nil {
			log.Println(err)
			return
		}
		if err := h.broker.SendMessage(msg); err != nil {
			// If there is an error it means that the broker is down
			// so we can just send the message to attached clients
			broadcast <- msg
			log.Printf("Err sendMessage: %v\n", err)
		}
	}
}
