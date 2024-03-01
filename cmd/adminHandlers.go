package cmd

import (
	"Platform/db"
	"context"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func getTeacherPage(w http.ResponseWriter, r *http.Request) {
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
	data := struct {
		Teachers []Teacher
	}{
		Teachers: teachers,
	}
	renderTemplate(w, "admin_pages/teachers.html", data)
}

func getStudentPage(w http.ResponseWriter, r *http.Request) {
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
		Students []Student
	}{
		Students: students,
	}
	renderTemplate(w, "admin_pages/students.html", data)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	collection := db.Client.Database("EduPortal").Collection("teachers")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Error deleting teacher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
}

func deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	collection := db.Client.Database("EduPortal").Collection("students")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Error deleting student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
}

func addStudentHandler(w http.ResponseWriter, r *http.Request) {
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

	newStudent := Student{
		Name:     name,
		Surname:  lastName,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	collection := db.Client.Database("EduPortal").Collection("students")
	_, err = collection.InsertOne(context.Background(), newStudent)
	if err != nil {
		http.Error(w, "Error saving student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/students", http.StatusSeeOther)
}

func AdminPanelHandler(w http.ResponseWriter, r *http.Request) {
	// fetching courses
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var courses []Courses
	coursesCollection := db.Client.Database("EduCourses").Collection("courses")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coursesCursor, err := coursesCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch courses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer coursesCursor.Close(ctx)

	for coursesCursor.Next(ctx) {
		var course Courses
		if err := coursesCursor.Decode(&course); err != nil {
			http.Error(w, "Failed to decode course: "+err.Error(), http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}
	if err := coursesCursor.Err(); err != nil {
		http.Error(w, "Cursor iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Courses []Courses
	}{
		Courses: courses,
	}

	if err := renderTemplate(w, "admin_pages/adminPanel.html", data); err != nil {
        http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
    }
}

func AdminAddCourseHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    name := r.FormValue("name")
    description := r.FormValue("description")
    category := r.FormValue("category")
    price, err := strconv.ParseFloat(r.FormValue("price"), 64)
    if err != nil {
        http.Error(w, "Invalid price", http.StatusBadRequest)
        return
    }
    imageUrl := r.FormValue("imageUrl")

    if name == "" || description == "" || category == "" {
        http.Error(w, "Fill the requirements", http.StatusBadRequest)
        return
    }

    newCourse := Courses {
        ID:          primitive.NewObjectID(),
        Name:        name,
        Description: description,
        Category:    category,
        Price:       price,
        Url:         imageUrl,
    }

    coursesCollection := db.Client.Database("EduCourses").Collection("courses")
    _, err = coursesCollection.InsertOne(context.Background(), newCourse)
    if err != nil {
        http.Error(w, "Failed to add new course: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/admin", http.StatusSeeOther)
}