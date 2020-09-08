package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.56.23:27017"))
	if err != nil {
		fmt.Printf("connect mongodb err: %v\n", err)
		return
	}
	if err := conn.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("Ping err: %v\n", err)
		return
	}

	conn.Database("my_db").Collection("my_collection")
	fmt.Println("pong")
	defer conn.Disconnect(ctx)
}
