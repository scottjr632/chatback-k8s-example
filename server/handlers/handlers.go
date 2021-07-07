package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/scottjr632/chatback/server/broker"
	"github.com/scottjr632/chatback/server/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func (h *Handlers) NewMessage(c *fiber.Ctx) error {
	msg := models.Message{}
	if err := c.BodyParser(&msg); err != nil {
		c.JSON(&fiber.Map{
			"status": false,
			"error":  err,
		})
		return err
	}

	if err := msg.Insert(context.TODO(), h.db, boil.Infer()); err != nil {
		log.Println(err)
		c.JSON(&fiber.Map{
			"status": false,
			"error":  err,
		})
		return err
	}

	broadcast <- msg
	c.JSON(&fiber.Map{
		"status": true,
		"error":  nil,
	})
	return nil
}
