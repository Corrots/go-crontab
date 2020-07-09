package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Log struct {
	// 任务名称
	TaskName string `bson:"task_name"`
	// shell指令
	Command string `bson:"command"`
	// shell执行error
	Error string `bson:"error"`
	// shell执行stdout
	Output string `bson:"output"`
	// 计划开始时间
	PlanTime string `bson:"plan_time"`
	// 实际调度时间
	ScheduleTime string `bson:"schedule_time"`
	// 任务执行开始时间
	StartTime string `bson:"start_time"`
	// 任务执行结束时间
	EndTime string `bson:"end_time"`
}

type Mongo struct {
	Ctx        context.Context
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongo() (*Mongo, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	uri := fmt.Sprintf("mongodb://%s", viper.GetString("mongo.addr"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongodb err: %v\n", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("Ping err: %v\n", err)
	}

	collection := client.Database("cron").Collection("log")
	return &Mongo{
		Ctx:        ctx,
		Client:     client,
		Collection: collection,
	}, nil
}

func (m *Mongo) InsertLog(log *Log) error {
	_, err := m.Collection.InsertOne(m.Ctx, log)
	defer m.Client.Disconnect(m.Ctx)
	return err
}
