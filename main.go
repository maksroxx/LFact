package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/maksroxx/LFact/api"
	"github.com/maksroxx/LFact/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	_ = client

	var (
		userStore = db.NewMongoUserStore(client)
		store     = &db.Store{
			UserStore: userStore,
		}
		userHandler = api.NewUserHandler(store)
		app         = fiber.New()
		apiv1       = app.Group("/api/v1")
	)

	// Versioned API routes
	// greeting
	apiv1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello lfact!")
	})

	// user handlers
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	if err := app.Listen(listenAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
