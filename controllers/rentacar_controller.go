// controllers/rentacar_controller.go

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"treci-mikroservis-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


var collection *mongo.Collection

func InitializeCollection(c *gin.Context) {
	
	client, exists := c.Get("mongoClient")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found"})
		return
	}

	collection = client.(*mongo.Client).Database("RentACar").Collection("UserCar")

 }

	func GetUserByID(userId int) (*models.User, error) {
		url := fmt.Sprintf("http://distribuirani-001-site1.btempurl.com/User/%d", userId)

		// Slanje HTTP GET zahteva
		response, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Failed to send HTTP request: %v", err)
		}
		defer response.Body.Close()

		// Provera status koda odgovora
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
		}

		// Čitanje i dekodiranje JSON odgovora
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read response body: %v", err)
		}

		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			return nil, fmt.Errorf("Failed to decode JSON response: %v", err)
		}

		return &user, nil
	}

func GetCarByID(carID int) (*models.Car, error) {
	url := fmt.Sprintf("http://sistemidis-001-site1.atempurl.com/Car/%d", carID)

	// Slanje HTTP GET zahteva
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to send HTTP request: %v", err)
	}
	defer response.Body.Close()

	// Provera status koda odgovora
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	// Čitanje i dekodiranje JSON odgovora
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	var car models.Car
	if err := json.Unmarshal(body, &car); err != nil {
		return nil, fmt.Errorf("Failed to decode JSON response: %v", err)
	}

	return &car, nil
}

func GetRentACar(c *gin.Context) {
	var rentACar []models.UserCar

	// Dobijanje MongoDB klijenta iz konteksta
	// client, exists := c.Get("mongoClient")
	// if !exists {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found"})
	// 	return
	// }

	// Dohvatanje MongoDB kolekcije
	//collection = client.(*mongo.Client).Database("RentACar").Collection("UserCar")

	// Kursor za dohvatanje rezultata
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from MongoDB"})
		return
	}
	defer cursor.Close(context.Background())

	// Iteriranje kroz rezultate i dodavanje u slice
	for cursor.Next(context.Background()) {
		var result models.UserCar


		// Dekodiranje rezultata
		if err := cursor.Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode data from MongoDB"})
			return
		}

		// Dodavanje u rezultate
		rentACar = append(rentACar, result)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while iterating through MongoDB results"})
		return
	}

	// Pravljenje novog niza koji sadrži samo id
	var response []gin.H
	for _, car := range rentACar {
		response = append(response, gin.H{"id": car.Id.Hex(), "userId": car.UserId, "carId": car.CarId,  "from": car.From, "to": car.To})
	}

	// Slanje JSON odgovora sa nizom koji sadrži samo id
	c.JSON(http.StatusOK, response)
}


func PostRentACar(c *gin.Context) {
	var rentACar models.UserCar
	var existingUserId models.UserCar

	print(rentACar.From)
	// Parsiranje JSON body-ja
	if err := c.ShouldBindJSON(&rentACar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Provera postojanja korisnika
	user, err := GetUserByID(rentACar.UserId)
	if err != nil {
		print(user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	car, err := GetCarByID(rentACar.CarId)
	if err != nil {
		print(car)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Car does not exist"})
		return
	}

	//proveri da li u collection userCar vec ima red sa userId
	collection.FindOne(context.Background(), bson.D{{"userid", rentACar.UserId}}).Decode(&existingUserId); 
	if existingUserId.UserId == rentACar.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already rented a car"})
		return
	}

	myDate, err := time.Parse("2006-01-02 15:04", rentACar.From)
	if err != nil {
		panic(err)
	}

	myDateTo, err := time.Parse("2006-01-02 15:04", rentACar.To)
	if err != nil {
		panic(err)
	}

	rentACar.From = myDate.Format("2006-01-02")
	rentACar.To = myDateTo.Format("2006-01-02")

	// Provera datuma iznajmljivanja
	if err := validateRentDates(rentACar, collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dodavanje dokumenta u kolekciju
	insertResult, err := collection.InsertOne(context.Background(), rentACar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into MongoDB"})
		return
	}

	// Slanje odgovora sa id-om dokumenta koji je dodat u kolekciju
	c.JSON(http.StatusOK, gin.H{"id": insertResult.InsertedID, "message": "You rented a car successfully"})
}
func GetRentACarById(c *gin.Context) {
	var rentACar models.UserCar

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id param"})
		return
	}

	// Dohvatanje dokumenta iz kolekcije
	if err := collection.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&rentACar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from MongoDB"})
		return
	}

	// Slanje JSON odgovora
	c.JSON(http.StatusOK, rentACar)
}

func DeleteRentACarById(c *gin.Context) {
	var rentACar models.UserCar

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id param"})
		return
	}

	// Brisanje dokumenta iz kolekcije
	if err := collection.FindOneAndDelete(context.Background(), bson.D{{"_id", id}}).Decode(&rentACar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data from MongoDB"})
		return
	}

	// Slanje JSON odgovora
	c.JSON(http.StatusOK, rentACar)
}


func validateRentDates(rentACar models.UserCar, collectionUserCar *mongo.Collection) error {

	
	// Parsiranje datuma iznajmljivanja u vrednosti vremena
	fromTime, err := time.Parse("2006-01-02", rentACar.From)
	if err != nil {
		return fmt.Errorf("Invalid date format for 'From'")
	}

	toTime, err := time.Parse("2006-01-02", rentACar.To)
	if err != nil {
		return fmt.Errorf("Invalid date format for 'To'")
	}

	// Provera datuma iznajmljivanja u kolekciji
	cursor, err := collectionUserCar.Find(context.Background(), bson.M{"carid": rentACar.CarId})
	if err != nil {
		return fmt.Errorf("Failed to query MongoDB for existing rentals")
	}
	defer cursor.Close(context.Background())

	// Iteracija kroz rezultate upita
	for cursor.Next(context.Background()) {
		var existingRent models.UserCar
		if err := cursor.Decode(&existingRent); err != nil {
			return fmt.Errorf("Failed to decode existing rental data")
		}

		existingFromTime, err := time.Parse("2006-01-02", existingRent.From)
		if err != nil {
			return fmt.Errorf("Invalid date format for 'From'")
		}

		existingToTime, err := time.Parse("2006-01-02", existingRent.To)
		if err != nil {
			return fmt.Errorf("Invalid date format for 'To'")
		}

		// Provera preklapanja vremenskih intervala
		if (fromTime.Before(existingToTime) && toTime.After(existingFromTime)) || (existingFromTime.Before(toTime) && existingToTime.After(fromTime)) {
			return fmt.Errorf("Car is already rented for the specified period")
		}
	}

	return nil
}
