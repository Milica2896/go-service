// main.go

package main

import (
	"treci-mikroservis-go/controllers"
	routes "treci-mikroservis-go/router"

	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Kreiranje Gin rutera
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://edinakucevic:39ejKp4t0k9AQTpE@cluster0.phbdhqv.mongodb.net/"))
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
			return
		}
		c.Set("mongoClient", client)
		c.Next()
	})

	router.Use(func(c *gin.Context) {
		controllers.InitializeCollection(c)
	})

	// Podešavanje ruta
	routes.SetupRentACarRoutes(router)

	// Pokretanje servera
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
