package core

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

func (m *Mongo) InsertMany(logs []interface{}) error {
	cxt := getContext()
	_, err := m.Collection.InsertMany(cxt, logs)
	//defer m.Client.Disconnect(cxt)
	return err
}
