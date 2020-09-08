package main

import (
	"context"
	"fmt"
	"time"

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

type earlyThan struct {
	LessThan int64 `bson:"$lt"`
}

type filter struct {
	EarlyThan earlyThan `bson:"timepoint.starttime"`
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
	filter := &filter{EarlyThan: earlyThan{time.Now().Unix()}}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		fmt.Printf("delete many err: %v\n", err)
		return
	}
	fmt.Printf("delete count: %d\n", result.DeletedCount)
}
