package db

import "os"

const MongoDBNameEnvName = "MONGO_DB_NAME"

var DBNAME = os.Getenv("MONGO_DB_URI")

type Store struct {
	UserStore UserStore
}

type Map map[string]any
