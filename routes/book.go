package routes

import (
	"bookstore/mongo-api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterBookRoutes enregistre les routes pour les opérations sur les livres
func RegisterBookRoutes(r *gin.Engine) {
	r.GET("/books", handlers.GetBooks)          // Récupérer tous les livres
	r.POST("/books", handlers.PostBooks)        // Ajouter un nouveau livre
	r.GET("/books/:id", handlers.GetBookByID)   // Récupérer un livre par ID
	r.PATCH("/books/:id", handlers.UpdateBook)  // Mettre à jour un livre par ID
	r.DELETE("/books/:id", handlers.DeleteBook) // Supprimer un livre par ID
}
