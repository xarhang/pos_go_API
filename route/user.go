package route

import (
	"pos-go/controller"
	"pos-go/middleware"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine, adapter *gormadapter.Adapter) {
	UserController := controller.User{}
	UserGroup := r.Group("/users",
		middleware.AuthorizationJWT(),
		middleware.Authorize("/users", "read", adapter),
	)
	UserGroup.GET("", UserController.FindAll)
	UserGroup.GET("/:id", UserController.FindOne)

}
