package cmd

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Surname  string             `json:"lastName" bson:"lastName"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password" bson:"password"`
	Student    *Student          `json:"student,omitempty" bson:"student,omitempty"`
}

type Student struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	Teacher *Teacher         	 `json:"teacher,omitempty" bson:"teacher,omitempty"`
}
