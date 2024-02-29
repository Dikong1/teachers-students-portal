package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func DbConnection() error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://azatabdirashituly:20feb2005@cluster0.sellypx.mongodb.net/")
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	CreateIndexes()

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")
	return nil
}
