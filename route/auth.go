package route

import (
	"pos-go/controller"

	"github.com/gin-gonic/gin"
)

func AuthorizationRoute(r *gin.Engine) {
	authController := controller.Auth{}
	authGroup := r.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.PATCH("/signin", authController.Signin)
}
