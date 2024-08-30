// handlers/book.go

package handlers

import (
	"context"
	"net/http"

	"bookstore/mongo-api/database"
	"bookstore/mongo-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetBooks retrieves all books from the MongoDB collection
func GetBooks(c *gin.Context) {
	var books []models.Book
	collection := database.GetClient().Database("bookstore").Collection("books")

	// Find all books
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	defer cursor.Close(context.Background())

	// Decode all books into the books slice
	if err = cursor.All(context.Background(), &books); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// PostBooks inserts a new book into the MongoDB collection
func PostBooks(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.GetClient().Database("bookstore").Collection("books")

	// Insert the new book
	result, err := collection.InsertOne(context.Background(), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insertedID": result.InsertedID})
}

// GetBookByID retrieves a book by its ID from the MongoDB collection
func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var book models.Book
	collection := database.GetClient().Database("bookstore").Collection("books")

	// Find the book by ID
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&book)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// UpdateBook updates an existing book by its ID
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.GetClient().Database("bookstore").Collection("books")

	// Define the update
	update := bson.M{
		"$set": bson.M{
			"title":  updatedBook.Title,
			"author": updatedBook.Author,
			"price":  updatedBook.Price,
		},
	}

	// Perform the update
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated"})
}

// DeleteBook deletes a book by its ID
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := database.GetClient().Database("bookstore").Collection("books")

	// Perform the delete
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting book"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
