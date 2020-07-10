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

type Mongo struct {
	Ctx        context.Context
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongo() (*Mongo, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(viper.GetInt("mongo.timeout")))
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
		Ctx:        ctx,
		Client:     client,
		Collection: collection,
	}, nil
}

//func (m *Mongo) InsertOne(log *core.Log) error {
//	_, err := m.Collection.InsertOne(m.Ctx, log)
//	defer m.Client.Disconnect(m.Ctx)
//	return err
//}
//
//func (m *Mongo) InsertMany(logs []interface{}) error {
//	_, err := m.Collection.InsertMany(m.Ctx, logs)
//	defer m.Client.Disconnect(m.Ctx)
//	return err
//}
