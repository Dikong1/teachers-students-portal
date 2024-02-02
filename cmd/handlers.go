package cmd

import (
	"Platform/db"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var log = logrus.New()

type errorss struct {
	ErrorCode int
	ErrorMsg  string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	err := renderTemplate(w, "home.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func teachLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		collection := db.Client.Database("EduPortal").Collection("teachers")
		var teacher Teacher
		if err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&teacher); err != nil {
			// Log and handle the error
			log.WithField("error", err).Error("Error finding teacher")
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password)); err != nil {
			// Log and handle the error
			log.WithField("error", err).Error("Password comparison failed")
			http.Error(w, "Invalid password: Password comparison failed", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/teach/%s", teacher.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "teachlog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func teachRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		teacher := Teacher{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
		}

		collection := db.Client.Database("EduPortal").Collection("teachers")

		result, err := collection.InsertOne(context.Background(), teacher)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/teach/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "teachreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func studLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		collection := db.Client.Database("EduPortal").Collection("students")
		var student Student
		if err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&student); err != nil {
			// Log and handle the error
			log.WithField("error", err).Error("Error finding student")
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password)); err != nil {
			// Log and handle the error
			log.WithField("error", err).Error("Password comparison failed")
			http.Error(w, "Invalid password: Password comparison failed", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/stud/%s", student.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "studlog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func studRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		student := Student{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
		}

		collection := db.Client.Database("EduPortal").Collection("students")

		result, err := collection.InsertOne(context.Background(), student)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		http.Redirect(w, r, fmt.Sprintf("/stud/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "studreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("frontend/templates/" + tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func teachPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["id"]

	var teacher Teacher
	collection := db.Client.Database("EduPortal").Collection("teachers")
	objID, _ := primitive.ObjectIDFromHex(teacherID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&teacher)
	if err != nil {
		http.Error(w, "teacher not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "teach.html", teacher)
}

func studPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	var student Student
	collection := db.Client.Database("EduPortal").Collection("students")
	objID, _ := primitive.ObjectIDFromHex(studentID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&student)
	if err != nil {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "stud.html", student)
}

func getDataFromDatabase(w http.ResponseWriter, r *http.Request) {
	collection := db.Client.Database("EduCourses").Collection("courses")

	filterValue := r.URL.Query().Get("filter")
	sortValue := r.URL.Query().Get("sort")

	var filter bson.M
	if filterValue != "" {
		filter = bson.M{"name": bson.M{"$regex": filterValue, "$options": "i"}}
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

func ErrorHandler(w http.ResponseWriter, r *http.Request, errCode int, msg string) {
	t, err := template.ParseFiles("frontend/templates/Error.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Errors := errorss{
		ErrorCode: errCode,
		ErrorMsg:  msg,
	}
	t.Execute(w, Errors)
}
