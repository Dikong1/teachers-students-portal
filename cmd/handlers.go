package cmd

import (
	"Platform/db"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var log = logrus.New()

var jwtKey = []byte("jr4fpiKTdbWFaVbXa1fs0mpI20MoJDTU")

type contextKey int

const (
	contextKeyUserID contextKey = iota
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log_.WithFields(logrus.Fields{
			"action": "homeHandler",
			"method": r.Method,
		}).Error("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		log_.WithFields(logrus.Fields{
			"action": "homeHandler",
			"path":   r.URL.Path,
		}).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := renderTemplate(w, "home.html", nil)
	if err != nil {
		log_.WithFields(logrus.Fields{
			"action": "homeHandler",
			"error":  err,
		}).Error("Internal server error")

		w.WriteHeader(http.StatusInternalServerError)
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
			log.WithField("error", err).Error("Error finding teacher")
			http.Error(w, "Cannot find email", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password)); err != nil {
			log.WithField("error", err).Error("Password comparison failed")
			http.Error(w, "Invalid password: Password comparison failed", http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID: teacher.ID.Hex(),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			http.Error(w, "Error signing token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		http.Redirect(w, r, fmt.Sprintf("/teach/%s", teacher.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "teachlog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
			log.WithField("error", err).Error("Error generating password hash")
			http.Error(w, "Error generating password hash", http.StatusInternalServerError)
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
			log.WithField("error", err).Error("Error inserting teacher")
			http.Error(w, "Error inserting teacher", http.StatusInternalServerError)
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/teach/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "teachreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
			log.WithField("error", err).Error("Error finding student")
			http.Error(w, "Cannot find email", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password)); err != nil {
			log.WithField("error", err).Error("Password comparison failed")
			http.Error(w, "Invalid password: Password comparison failed", http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(1 * time.Hour) // Token valid for 1 hour
		claims := &Claims{
			UserID: student.ID.Hex(),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.WithField("error", err).Error("Error signing token")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		http.Redirect(w, r, fmt.Sprintf("/stud/%s", student.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "studlog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

		collection := db.Client.Database("EduPortal").Collection("students")

		result, err := collection.InsertOne(context.Background(), student)
		if err != nil {
			log.WithField("error", err).Error("Error inserting stdent")
			http.Error(w, "Error inserting student", http.StatusInternalServerError)
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		http.Redirect(w, r, fmt.Sprintf("/stud/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "studreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	var students []Student
	studentCollection := db.Client.Database("EduPortal").Collection("students")
	cursor, err := studentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.WithField("error", err).Error("Error fetching students")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var student Student
		if err = cursor.Decode(&student); err != nil {
			log.WithField("error", err).Error("Error decoding student data")
			continue
		}
		students = append(students, student)
	}

	data := struct {
		Teacher  Teacher
		Students []Student
	}{
		Teacher:  teacher,
		Students: students,
	}

	renderTemplate(w, "teach.html", data)
}

func studPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	// Log the student ID being accessed
	log.WithField("studentID", studentID).Info("Accessing student page")

	collection := db.Client.Database("EduPortal").Collection("students")
	objID, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		log.WithField("error", err).Error("Invalid student ID format")
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	var student Student
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.WithField("studentID", studentID).Warn("No student found with given ID")
			http.Error(w, "No student found", http.StatusNotFound)
		} else {
			log.WithField("error", err).Error("Error finding student")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Log successful data retrieval
	log.WithField("student", student).Info("Successfully retrieved student data")

	err = renderTemplate(w, "stud.html", student)
	if err != nil {
		log.WithField("error", err).Error("Error rendering template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func verifyToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
