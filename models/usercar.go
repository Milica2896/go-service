package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

)
type UserCar struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId int                `json:"userId"`
	CarId  int                `json:"carId"`
	From   string 			  `json:"from"`
	To     string 			  `json:"to"`
}
