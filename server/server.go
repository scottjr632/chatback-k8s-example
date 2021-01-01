//go:generate sqlboiler psql
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/scottjr632/chatback/server/config"
	"github.com/scottjr632/chatback/server/handlers"

	"github.com/gofiber/websocket/v2"
	_ "github.com/lib/pq"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func init() {
	godotenv.Load()
}

func createDB(dbConfig config.DB) (*sql.DB, error) {
	dbString := fmt.Sprintf("dbname=%s user=%s sslmode=disable password=%s host=%s",
		dbConfig.DBName, dbConfig.User, dbConfig.Password, dbConfig.Host)
	return sql.Open("postgres", dbString)
}

func main() {
	config := config.New()
	config.Broker.URI = os.Getenv("BROKER_URI")
	db, err := createDB(*config.DB)
	if err != nil {
		panic(err)
	}

	h := handlers.New(db, config.Broker)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Get("/api/messages", h.GetMessages)
	app.Use("/ws", h.UpgradeWebsocket)
	app.Get("/ws/messages", websocket.New(h.WebsocketMessages))

	log.Fatal(app.Listen(":8080"))
}
