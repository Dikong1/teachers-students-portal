package cmd

import (
	"Platform/db"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

)

func getDataFromDatabase(w http.ResponseWriter, r *http.Request) {
	collection := db.Client.Database("EduCourses").Collection("courses")

	filterValue := r.URL.Query().Get("filter")
	sortValue := r.URL.Query().Get("sort")

	var filter bson.M
	if filterValue != "" {
		filter = bson.M{"category": bson.M{"$regex": filterValue, "$options": "i"}}
	} else {
		filter = bson.M{}
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	const pageSize = 6
	skip := int64((page - 1) * pageSize)

	findOptions := options.Find()
	if sortValue != "" {
		switch sortValue {
		case "name":
			findOptions.SetSort(bson.D{{Key: "name", Value: 1}})
		case "price":
			findOptions.SetSort(bson.D{{Key: "price", Value: 1}})
		}
	}

	findOptions.SetSkip(skip)
	findOptions.SetLimit(pageSize)

	// LOGRUS OPERATION
	log.WithFields(logrus.Fields{
		"action": "fetch_data",
		"filter": filterValue,
		"sort":   sortValue,
		"page":   page,
	}).Info("Fetching data from database")

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		//LOGRUS OPERATION
		log.WithFields(logrus.Fields{
			"action": "fetch_error",
			"error":  err.Error(),
		}).Error("Error fetching data from database")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []Courses
	for cursor.Next(context.Background()) {
		var course Courses
		if err := cursor.Decode(&course); err != nil {
			//LOGRUS OPERATION
			log.WithFields(logrus.Fields{
				"action": "decode_error",
				"error":  err.Error(),
			}).Error("Error decoding course data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	if err := cursor.Err(); err != nil {
		//LOGRUS OPERATION
		log.WithFields(logrus.Fields{
			"action": "json_encode_error",
			"error":  err.Error(),
		}).Error("Error encoding courses to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(courses)
	if err != nil {
		//LOGRUS OPERATION
		log.WithFields(logrus.Fields{
			"action": "json_encode_error",
			"error":  err.Error(),
		}).Error("Error encoding courses to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	log.WithFields(logrus.Fields{
		"action":     "data_fetched",
		"numCourses": len(courses),
	}).Info("Successfully fetched course data")
}
