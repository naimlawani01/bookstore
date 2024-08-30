package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"` // ID de type ObjectID pour MongoDB
	Title  string             `json:"title" bson:"title"`
	Author string             `json:"author" bson:"author"`
	Price  float64            `json:"price" bson:"price"`
}
