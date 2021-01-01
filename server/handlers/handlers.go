package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/scottjr632/chatback/server/broker"
	"github.com/scottjr632/chatback/server/models"
)

type Handlers struct {
	db     *sql.DB
	broker *broker.Broker
}

func New(db *sql.DB, bc *broker.Config) *Handlers {
	b := broker.New(bc)
	h := &Handlers{db, b}
	go h.handleRegisters()
	return &Handlers{db, b}
}

func (h *Handlers) GetMessages(c *fiber.Ctx) error {
	messages, err := models.Messages().All(c.Context(), h.db)
	if err != nil {
		return err
	}
	c.JSON(messages)
	return nil
}
