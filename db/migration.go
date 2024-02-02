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
			Keys:    map[string]interface{}{"phone": 1},
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
			Keys:    map[string]interface{}{"phone": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	// courseCollection := client.Database("EduPortal").Collection("courses")
	// _, err = courseCollection.Indexes().CreateOne(
	// 	context.Background(),
	// 	mongo.IndexModel{
	// 		Keys:    map[string]interface{}{"name": 1},
	// 		Options: options.Index().SetUnique(true),
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }
	// filePath := filepath.Join("./courses.json")
	// file, err := os.ReadFile(filePath)
	// if err != nil {
	// 	log.Fatalf("Error reading JSON file: %v", err)
	// }

	// var courses []cmd.Courses

	// if err := json.Unmarshal(file, &courses); err != nil {
	// 	log.Fatalf("Error unmarshaling JSON data: %v", err)
	// }

	// var docs []interface{}
	// for _, course := range courses {
	// 	docs = append(docs, course)
	// }
	// course_collection := client.Database("EduPortal").Collection("courses")
	// _, err = course_collection.InsertMany(context.Background(), docs)
	// if err != nil {
	// 	log.Fatalf("Error inserting data into MongoDB: %v", err)
	// }

	// log.Println("Courses inserted successfully")

	return nil
}
