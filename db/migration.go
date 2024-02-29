package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() error {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	teachersCollection := client.Database("EduPortal").Collection("teachers")
	_, err = teachersCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	studentsCollection := client.Database("EduPortal").Collection("students")
	_, err = studentsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
