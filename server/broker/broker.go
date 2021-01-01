package broker

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/r3labs/sse/v2"
	"github.com/scottjr632/chatback/server/models"
)

type Broker struct {
	*Config
}

func New(c *Config) *Broker {
	return &Broker{c}
}

func (b *Broker) SendMessage(message models.Message) error {
	client := http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	event := &sse.Event{
		ID:   []byte(strconv.FormatInt(int64(message.ID), 10)),
		Data: messageJSON,
	}

	eventBytes, err := json.Marshal(event)
	reader := bytes.NewReader(eventBytes)

	resp, err := client.Post(b.URI+"/message", "application/json", reader)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New("Invalid status code")
	}

	return nil
}

// HandleSubscribe ...
func (b *Broker) HandleSubscribe() (chan *sse.Event, chan error) {
	client := sse.NewClient(b.URI + "/subscribe")
	client.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	ch := make(chan *sse.Event)
	errChan := make(chan error)
	go func() {
		errChan <- client.SubscribeChan("messages", ch)
	}()
	return ch, errChan
}
