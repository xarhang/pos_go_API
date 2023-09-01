package route

import (
	"pos-go/controller"
	"pos-go/middleware"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

func PermissionRoute(r *gin.Engine, adapter *gormadapter.Adapter) {
	PermissionController := controller.Permission{}
	PermissionGroup := r.Group("/permissions",
		middleware.AuthorizationJWT(),
	)
	//get all permission
	PermissionGroup.GET("",
		middleware.Authorize("/permissions", "read", adapter),
		PermissionController.FindAll,
	)
	PermissionGroup.GET("/:rule",
		middleware.Authorize("/permissions", "read", adapter),
		PermissionController.FindOne,
	)
	//create new permission
	PermissionGroup.POST("",
		middleware.Authorize("/permissions", "write", adapter),
		PermissionController.Create,
	)
	//Delete permission
	PermissionGroup.DELETE("",
		middleware.Authorize("/permissions", "write", adapter),
		PermissionController.Delete,
	)

}
