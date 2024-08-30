package main

import (
	"bookstore/mongo-api/database"
	"bookstore/mongo-api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Charger les variables d'environnement
	if err := database.ConnectMongoDB("mongodb+srv://najathlawani73:QbJJXgG06MhCg02z@cluster0.4kgjy6f.mongodb.net/"); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.DisconnectMongoDB() // Ensure the connection is closed when the application exits
	router := gin.Default()
	routes.RegisterBookRoutes(router)
	router.Run("localhost:8080")
}
