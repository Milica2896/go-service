
package models


import ()

type Car struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImagePath string `json:"imagePath"`
	Price int `json:"price"`
	Year int `json:"year"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Fuel string `json:"fuel"`
	Transmission string `json:"transmission"`
}