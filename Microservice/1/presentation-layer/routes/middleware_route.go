package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/middleware"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/presentation-layer/controller"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("/api/user-management")
	{
		routes.GET("/user-by-token", middleware.Authenticate(jwtService), userController.GetUserByToken)
	}
}
