package cmd

import (
	"Platform/db"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func getAddingTeacherPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "admin_pages/addTeacher.html", nil)
}

func getAddingStudentPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "admin_pages/addStudent.html", nil)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// First, parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	newTeacher := Teacher{
		Name:     name,
		Surname:  lastName,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	collection := db.Client.Database("EduPortal").Collection("teachers")
	_, err = collection.InsertOne(context.Background(), newTeacher)
	if err != nil {
		http.Error(w, "Error saving teacher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []Teacher
	teacherCollection := db.Client.Database("EduPortal").Collection("teachers")
	teacherCursor, err := teacherCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch teachers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer teacherCursor.Close(context.Background())

	for teacherCursor.Next(context.Background()) {
		var teacher Teacher
		if err := teacherCursor.Decode(&teacher); err != nil {
			http.Error(w, "Failed to decode teacher: "+err.Error(), http.StatusInternalServerError)
			return
		}
		teachers = append(teachers, teacher)
	}

	// Fetch all students
	var students []Student
	studentCollection := db.Client.Database("EduPortal").Collection("students")
	studentCursor, err := studentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch students: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer studentCursor.Close(context.Background())

	for studentCursor.Next(context.Background()) {
		var student Student
		if err := studentCursor.Decode(&student); err != nil {
			http.Error(w, "Failed to decode student: "+err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	data := struct {
		Teachers []Teacher
		Students []Student
	}{
		Teachers: teachers,
		Students: students,
	}

	renderTemplate(w, "admin_pages/adminPanel.html", data)
}
