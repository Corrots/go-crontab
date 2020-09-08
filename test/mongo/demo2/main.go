package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type LogRecord struct {
	JobName   string `json:"job_name"`
	Command   string `json:"command"`
	Err       string `json:"err"`
	Content   string `json:"content"`
	TimePoint TimePoint
}

type TimePoint struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

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

	collection := conn.Database("cron").Collection("log")
	//
	now := time.Now().Unix()
	record := &LogRecord{
		JobName: "job10",
		Command: "sleep 5;echo hello",
		Err:     "",
		Content: "The Last Of Us",
		TimePoint: TimePoint{
			StartTime: now,
			EndTime:   now + rand.Int63n(999),
		},
	}
	res, err := collection.InsertOne(ctx, record)
	if err != nil {
		fmt.Printf("insert err: %v\n", err)
		return
	}
	fmt.Printf("InsertedID: %b\n", res.InsertedID.(primitive.ObjectID))
	defer conn.Disconnect(ctx)
}
