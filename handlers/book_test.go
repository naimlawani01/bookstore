package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookstore/mongo-api/database"
	"bookstore/mongo-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const mongoURI = "mongodb+srv://najathlawani73:QbJJXgG06MhCg02z@cluster0.4kgjy6f.mongodb.net/"

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", GetBooks)
	r.POST("/books", PostBooks)
	r.GET("/books/:id", GetBookByID)
	return r
}

func setupDatabase() {
	err := database.ConnectMongoDB(mongoURI)
	if err != nil {
		panic(err)
	}
}

func teardownDatabase() {
	database.DisconnectMongoDB()
}

func clearCollection() {
	collection := database.GetClient().Database("bookstore_test").Collection("books")
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
}

func TestGetBooks(t *testing.T) {
	setupDatabase()
	defer teardownDatabase()

	clearCollection()

	r := setupRouter()

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostBooks(t *testing.T) {
	setupDatabase()
	defer teardownDatabase()

	clearCollection()
	collection := database.GetClient().Database("bookstore").Collection("books")
	r := setupRouter()

	book := models.Book{
		Title:  "New Book",
		Author: "New Author",
		Price:  29.95,
	}

	bookJSON, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ------- copy and past from search ---------
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	// ---------------
	_, ok := response["insertedID"]
	insertedID, ok := response["insertedID"].(string)
	assert.True(t, ok)
	collection.DeleteOne(context.Background(), bson.M{"_id": insertedID})
}

func TestGetBookByID(t *testing.T) {
	setupDatabase()
	defer teardownDatabase()

	clearCollection()

	collection := database.GetClient().Database("bookstore").Collection("books")
	testBook := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Price:  19.99,
	}
	insertResult, err := collection.InsertOne(context.Background(), testBook)
	if err != nil {
		t.Fatalf("Failed to insert test book: %v", err)
	}

	bookID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	r := setupRouter()
	req, _ := http.NewRequest("GET", "/books/"+bookID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	collection.DeleteOne(context.Background(), bson.M{"_id": bookID})

}
