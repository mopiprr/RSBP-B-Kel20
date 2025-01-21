package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/middleware"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/presentation-layer/controller"
)

func Log(route *gin.Engine, logController controller.LogController, jwtService service.JWTService) {
	routes := route.Group("/api/log")
	{
		routes.POST("/create-log", middleware.Authenticate(jwtService), logController.CreateLog)
		routes.GET("/get-log-user", middleware.Authenticate(jwtService), logController.GetLogByUser)
	}
}
