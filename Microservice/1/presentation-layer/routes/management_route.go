package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/middleware"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/presentation-layer/controller"
)

func Management(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("/api/management")
	{
		routes.GET("/users", middleware.Authenticate(jwtService), userController.GetAllUser)
	}
}
