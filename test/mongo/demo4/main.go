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

type filter struct {
	JobName string `json:"job_name"`
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
	filter := &filter{
		JobName: "job10",
	}
	var ops options.FindOptions
	cursor, err := collection.Find(ctx, filter, ops.SetSkip(0).SetLimit(4))
	if err != nil {
		fmt.Printf("find err: %v\n", err)
		return
	}
	var rows []LogRecord
	for cursor.Next(ctx) {
		var row LogRecord
		if err := cursor.Decode(&row); err != nil {
			fmt.Printf("decode record err: %v\n", err)
			break
		}
		rows = append(rows, row)
	}
	// 释放 cursor
	defer cursor.Close(ctx)
	fmt.Println(rows)
}
