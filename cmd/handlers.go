package cmd

import (
	"Platform/db"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var log = logrus.New()

var jwtKey = []byte("jr4fpiKTdbWFaVbXa1fs0mpI20MoJDTU")

type contextKey int

const (
	contextKeyUserID contextKey = iota
)

// var (
// 	mgDomain = os.Getenv("MG_DOMAIN")
// 	mgAPIKey = os.Getenv("MG_API_KEY")
// )

var (
	verificationTokens = make(map[string]string)
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
		Log.WithFields(logrus.Fields{
			"action": "homeHandler",
			"method": r.Method,
		}).Error("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		Log.WithFields(logrus.Fields{
			"action": "homeHandler",
			"path":   r.URL.Path,
		}).Error("Not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := renderTemplate(w, "home.html", nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
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
			Log.WithField("error", err).Error("Error finding teacher")
			http.Error(w, "Cannot find email", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password)); err != nil {
			Log.WithField("error", err).Error("Password comparison failed")
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

var tempTeacher Teacher // Temporary global variable to store teacher data

func teachRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		// Generate a random verification token
		verificationToken, err := generateRandomToken()
		if err != nil {
			Log.WithField("error", err).Error("Error generating verification token")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Store the verification token along with the user's email address
		verificationTokens[email] = verificationToken

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			Log.WithField("error", err).Error("Error generating password hash")
			http.Error(w, "Error generating password hash", http.StatusInternalServerError)
			return
		}

		// Store the teacher data in the temporary global variable
		tempTeacher = Teacher{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
		}

		// Send verification email
		err = sendVerificationEmail(email, verificationToken, "teacher")
		if err != nil {
			Log.WithField("error", err).Error("Error sending verification email")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Render a page indicating that the verification email has been sent
		renderTemplate(w, "verification_pending.html", nil)
		return
	} else if r.Method == "GET" {
		renderTemplate(w, "teachreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the token and email from the query parameters
	token := r.URL.Query().Get("token")
	email := r.URL.Query().Get("email")
	who := r.URL.Query().Get("who")

	// Perform verification using the handleVerification function
	handleVerification(w, r, token, email, who)
}

func verifyFailureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "verification_failure.html", nil)
	}
}

func sendVerificationEmail(email, token string, who string) error {
	// Create a new Mailgun client with your domain and API key
	mg := mailgun.NewMailgun("your-mailgun-domain", "your-mailgun-api-key")

	// Compose the email body with the verification link containing the token
	body := fmt.Sprintf("Click the link below to verify your email:\n\nhttp://localhost:3000/verify?token=%s&email=%s&who=%s", token, email, who)

	// Compose the email message
	message := mg.NewMessage(
		"Excited User <mailgun@your-domain.com>", // Replace with sender email address
		"Email Verification",                     // Email subject
		body,                                     // Email body
		email,                                    // Recipient email address
	)

	// Send the email using the Mailgun client
	_, _, err := mg.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}

func handleVerification(w http.ResponseWriter, r *http.Request, token, email string, who string) {
	// Check if the token is valid
	if storedToken, ok := verificationTokens[email]; ok && storedToken == token {
		// Token is valid, store the teacher data in the database

		if who == "teacher" {
			err := storeTeacherData(tempTeacher)
			if err != nil {
				Log.WithField("error", err).Error("Error storing teacher data")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else if who == "student" {
			err := storeStudentData(tempStudent)
			if err != nil {
				Log.WithField("error", err).Error("Error storing student data")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Clear the temporary global variable
		tempTeacher = Teacher{}
		tempStudent = Student{}

		// Redirect the user to a page indicating successful verification
		if who == "teacher" {
			http.Redirect(w, r, "/teachlog", http.StatusSeeOther)
		} else if who == "student" {
			http.Redirect(w, r, "/studlog", http.StatusSeeOther)
		}

		return
	}

	// Token is invalid or expired
	// Redirect the user to a page indicating failed verification
	http.Redirect(w, r, "/verification-failure", http.StatusSeeOther)
}

func storeTeacherData(teacher Teacher) error {
	collection := db.Client.Database("EduPortal").Collection("teachers")

	result, err := collection.InsertOne(context.Background(), teacher)
	if err != nil {
		return err
	}

	Log.Info("Teacher data stored successfully:", result.InsertedID)
	return nil
}

func storeStudentData(student Student) error {
	collection := db.Client.Database("EduPortal").Collection("students")

	result, err := collection.InsertOne(context.Background(), student)
	if err != nil {
		return err
	}

	Log.Info("Student data stored successfully:", result.InsertedID)
	return nil
}

func generateRandomToken() (string, error) {
	// Define the length of the token (adjust as needed)
	tokenLength := 32

	// Create a byte slice to store the random token
	tokenBytes := make([]byte, tokenLength)

	// Read random bytes into the tokenBytes slice
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", errors.New("error generating random token: " + err.Error())
	}

	// Encode the random bytes to base64 to create a string token
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	fmt.Println(token)

	return token, nil
}

func studLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		collection := db.Client.Database("EduPortal").Collection("students")
		var student Student
		if err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&student); err != nil {
			Log.WithField("error", err).Error("Error finding student")
			http.Error(w, "Cannot find email", http.StatusBadRequest)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password)); err != nil {
			Log.WithField("error", err).Error("Password comparison failed")
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
			Log.WithField("error", err).Error("Error signing token")
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

var tempStudent Student

func studRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		// Generate a random verification token
		verificationToken, err := generateRandomToken()
		if err != nil {
			Log.WithField("error", err).Error("Error generating verification token")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Store the verification token along with the user's email address
		verificationTokens[email] = verificationToken

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			Log.WithField("error", err).Error("Error generating password hash")
			http.Error(w, "Error generating password hash", http.StatusInternalServerError)
			return
		}

		// Store the student data in the temporary global variable while verification
		tempStudent = Student{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
		}

		// Send verification email
		err = sendVerificationEmail(email, verificationToken, "student")
		if err != nil {
			Log.WithField("error", err).Error("Error sending verification email")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Render a page indicating that the verification email has been sent
		renderTemplate(w, "verification_pending.html", nil)
		return
	} else if r.Method == "GET" {
		renderTemplate(w, "studreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
		Log.WithField("error", err).Error("Error fetching students")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var student Student
		if err = cursor.Decode(&student); err != nil {
			Log.WithField("error", err).Error("Error decoding student data")
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
