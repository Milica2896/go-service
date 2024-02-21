// routes/rentacar_routes.go

package routes

import (
	"log"
	"treci-mikroservis-go/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRentACarRoutes(router *gin.Engine) {
	rentACarGroup := router.Group("/api/")

	// Postavljanje osnovnih ruta
	rentACarGroup.GET("usercars", controllers.GetRentACar)
	rentACarGroup.POST("create-usercar", controllers.PostRentACar)
	rentACarGroup.GET("usercar/:id", controllers.GetRentACarById)
	rentACarGroup.DELETE("delete-usercar/:id", controllers.DeleteRentACarById)

	log.Println("pozval fju IZ RUTERA:")
}
