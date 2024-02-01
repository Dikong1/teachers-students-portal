package cmd

import (
	"SantaWeb/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

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

func volLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		incMsg := "Wrong password or phone"

		collection := db.Client.Database("EduPortal").Collection("teachers")
		var volunteer Volunteer
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&volunteer)
		if err != nil {
			renderTemplate(w, "vollogin.html", incMsg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(volunteer.Password), []byte(password))
		if err != nil {
			renderTemplate(w, "vollogin.html", incMsg)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/vol/%s", volunteer.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "vollogin.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func volRegHandler(w http.ResponseWriter, r *http.Request) {
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

		volunteer := Volunteer{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
			Child:    nil,
		}

		collection := db.Client.Database("EduPortal").Collection("teachers")

		result, err := collection.InsertOne(context.Background(), volunteer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/vol/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "volreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func chiLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		incMsg := "Wrong password or phone"

		collection := db.Client.Database("EduPortal").Collection("students")
		var child Child
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&child)
		if err != nil {
			renderTemplate(w, "chilog.html", incMsg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(child.Password), []byte(password))
		if err != nil {
			renderTemplate(w, "chilog.html", incMsg)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/chil/%s", child.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "chilog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func chiRegHandler(w http.ResponseWriter, r *http.Request) {
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

		child := Child{
			Name:      firstName,
			Surname:   lastName,
			Email:     email,
			Phone:     phone,
			Password:  string(hashedPassword),
			Wish:      nil,
			Volunteer: nil,
		}

		collection := db.Client.Database("EduPortal").Collection("students")
		collectionWishes := db.Client.Database("SantaWeb").Collection("wishes")

		result, err := collection.InsertOne(context.Background(), child)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		wishes := Wish{
			Wishes: "",
			Child:  &child,
		}
		_, err = collectionWishes.InsertOne(context.Background(), wishes)
		if err != nil {
			http.Error(w, "Error creating wish", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/chil/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		renderTemplate(w, "chireg.html", nil)
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

func volunteerPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volunteerID := vars["id"]

	var volunteer Volunteer
	collection := db.Client.Database("EduPortal").Collection("teachers")
	objID, _ := primitive.ObjectIDFromHex(volunteerID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Volunteer not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "vol.html", volunteer)
}

func childPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childID := vars["id"]

	var child Child
	collection := db.Client.Database("EduPortal").Collection("students")
	objID, _ := primitive.ObjectIDFromHex(childID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&child)
	if err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "chil.html", child)
}

/*func updateWishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		childIdString := r.FormValue("childId")
		newWish := r.FormValue("wish")

		childId, err := primitive.ObjectIDFromHex(childIdString)
		if err != nil {
			http.Error(w, "Invalid Child ID", http.StatusBadRequest)
			return
		}

		collection := db.Client.Database("SantaWeb").Collection("Children")

		filter := bson.M{"_id": childId}
		update := bson.M{"$set": bson.M{"wish": newWish}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/path-after-updating-wish", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}*/

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func updateWishesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	childID, err := getChildIDFromSession(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	wishes := r.FormValue("wishes")

	wishesCollection := db.Client.Database("SantaWeb").Collection("wishes")

	filter := bson.M{"childId": childID}
	update := bson.M{"$set": bson.M{"wishes": wishes}}
	_, err = wishesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating wishes", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/path-to-child-page", http.StatusSeeOther)
}

func getChildIDFromSession(r *http.Request) (primitive.ObjectID, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return primitive.NilObjectID, err
	}

	childID, ok := session.Values["childID"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("Child ID not found in session")
	}

	objID, err := primitive.ObjectIDFromHex(childID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objID, nil
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, errCode int, msg string) {
	t, err := template.ParseFiles("frontend/templates/Error.html")
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Errors := errorss{
		ErrorCode: errCode,
		ErrorMsg:  msg,
	}
	// w.WriteHeader(Errors.ErrorCode)
	t.Execute(w, Errors)
}
