package main

import (
	"context"
	"log"
	"os"

	"time"

	"vault-test/db"
	"vault-test/vault"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func routine(tick *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <- tick.C :
			log.Print("check connection to mongo")
			db.Reconnect()
			break;
		case <- quit :
			log.Print("stop routine")
			return;
		}
	}
}

func main() {
	godotenv.Load(".env")
	vault.Login()
	vault.GetDatabaseCred( os.Getenv("VAULT_TOKEN") )
	tick := time.NewTicker( time.Minute * 1 )
	quit := make(chan struct{})

	go routine(tick, quit)
	
	db.Connect()
	defer db.DB().Disconnect(context.TODO())

	time.Sleep(time.Second * 30)
	res := db.DB().Database("chatapp").Collection("users").FindOne(context.TODO(), bson.D{})
	
	if err := res.Err(); err!=nil {
		print(err.Error())
		db.Reconnect()
	}

	time.Sleep(time.Minute * 5)
	close(quit)
}