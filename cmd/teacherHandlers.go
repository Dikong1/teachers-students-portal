package cmd

import (
	"Platform/db"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func addStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.WithField("error", err).Error("Error parsing form")
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("name")
	lastName := r.FormValue("surname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithField("error", err).Error("Error generating password hash")
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	student := Student{
		Name:     firstName,
		Surname:  lastName,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	studentCollection := db.Client.Database("EduPortal").Collection("students")
	count, err := studentCollection.CountDocuments(context.Background(), student)
	if err != nil {
		log.WithField("error", err).Error("Error checking student")
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Student with the same email or phone already exists", http.StatusBadRequest)
		return
	}

	_, err = studentCollection.InsertOne(context.Background(), student)
	if err != nil {
		log.WithField("error", err).Error("Error adding student")
		http.Error(w, "Error adding student", http.StatusInternalServerError)
		return
	}

	teacherID := r.FormValue("teacherID")
	http.Redirect(w, r, fmt.Sprintf("/teach/%s", teacherID), http.StatusSeeOther)
}


func deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.WithField("error", err).Error("Error parsing form")
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }
    studentID := r.FormValue("studentID")

    objID, err := primitive.ObjectIDFromHex(studentID)
    if err != nil {
        log.WithField("error", err).Error("Invalid student ID")
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    _, err = db.Client.Database("EduPortal").Collection("students").DeleteOne(context.Background(), bson.M{"_id": objID})
    if err != nil {
        log.WithField("error", err).Error("Error deleting student")
        http.Error(w, "Error deleting student", http.StatusInternalServerError)
        return
    }

    teacherID := r.FormValue("teacherID")
    http.Redirect(w, r, fmt.Sprintf("/teach/%s", teacherID), http.StatusSeeOther) 
}
