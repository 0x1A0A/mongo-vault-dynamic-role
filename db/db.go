package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"vault-test/vault"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Client

func DB() *mongo.Client {
	return db
}

func Connect() {
	db = nil
	server, ok := os.LookupEnv("DB_SERVER")
	if !ok {
		server = "0.0.0.0"
	}
	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		port = "27017"
	}
	user, _ := os.LookupEnv("DB_USER")
	passwd, _ := os.LookupEnv("DB_PASSWD")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/chatapp", user, passwd, server, port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err!=nil {
		log.Print(err.Error())
		return
	}

	db = client
	if Ping() {
		log.Print("mongodb connected!")
	} else {
		log.Fatal("Cannot connect to mongodb")
	}
}

func Ping() bool {
	if db==nil {
		return false
	}
	_, err := db.ListDatabaseNames(context.TODO(), bson.D{{}})
	if err != nil {
		log.Print(err.Error())
		return false
	} else {
		return true
	}
}

func Reconnect() {
	if !Ping() { // need reconnect
		log.Print("reconnecting to mongo!")
		if err := db.Disconnect(context.TODO()); err!=nil {
			log.Print("can't disconnect old database conection")
			return
		}
		vault.Login()
		vault.GetDatabaseCred(os.Getenv("VAULT_TOKEN"))
		Connect()
	} else {
		log.Print("Ok!")
	}
}