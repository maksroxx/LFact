package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/maksroxx/LFact/db"
	"github.com/maksroxx/LFact/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URI")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{UserStore: db.NewMongoUserStore(client)}

	// add users
	fixtures.AddUser(store, "James", "Foo")
	fixtures.AddUser(store, "Bob", "Johnson")
	fixtures.AddUser(store, "Mark", "from the block")
}
