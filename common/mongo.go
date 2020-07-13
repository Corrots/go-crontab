package common

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongo() (*Mongo, error) {
	ctx := context.Background()
	uri := fmt.Sprintf("mongodb://%s", viper.GetString("mongo.addr"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongodb err: %v\n", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("Ping err: %v\n", err)
	}

	collection := client.Database(viper.GetString("mongo.database")).Collection(viper.GetString("mongo.collection"))
	return &Mongo{
		Client:     client,
		Collection: collection,
	}, nil
}

func getContext() context.Context {
	return context.Background()
}

func (m *Mongo) InsertLogs(logs []interface{}) error {
	_, err := m.Collection.InsertMany(context.TODO(), logs)
	//defer m.Client.Disconnect(cxt)
	return err
}

type NameFilter struct {
	TaskName string `bson:"taskname"`
}

type LogSort struct {
	Sort int `bson:"starttime"`
}

func (m *Mongo) GetLogs(taskName string, skip, limit int64) ([]Log, error) {
	ctx := getContext()
	var ops options.FindOptions
	sort := LogSort{Sort: -1}
	cursor, err := m.Collection.Find(ctx, NameFilter{TaskName: taskName}, ops.SetSkip(skip).SetLimit(limit).SetSort(sort))
	if err != nil {
		return nil, err
	}
	var rows []Log
	for cursor.Next(ctx) {
		var row Log
		if err := cursor.Decode(&row); err != nil {
			return nil, fmt.Errorf("cursor log decode err: %v\n", err)
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no logs found")
	}
	defer cursor.Close(ctx)
	return rows, nil
}
