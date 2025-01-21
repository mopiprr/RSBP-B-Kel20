package main

import (
	"log"
	"os"

	"github.com/mci-its/backend-service/data-layer/config"
	"github.com/mci-its/backend-service/data-layer/repository"
	"github.com/mci-its/backend-service/domain-layer/middleware"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/helpers/cmd"
	"github.com/mci-its/backend-service/presentation-layer/controller"
	"github.com/mci-its/backend-service/presentation-layer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if len(os.Args) > 1 {
		cmd.Commands(db)
		return
	}

	var (
		jwtService service.JWTService = service.NewJWTService()

		// Implementation Dependency Injection
		// Repository
		userRepository   repository.UserRepository   = repository.NewUserRepository(db)
		roleRepository   repository.RoleRepository   = repository.NewRoleRepository(db)
		logRepository    repository.LogRepository    = repository.NewLogRepository(db)
		officeRepository repository.OfficeRepository = repository.NewOfficeRepository(db)

		// Service
		userService service.UserService = service.NewUserService(userRepository, roleRepository, officeRepository, jwtService)
		logService  service.LogService  = service.NewLogService(logRepository, userRepository)

		// Controller
		userController controller.UserController = controller.NewUserController(userService)
		logController  controller.LogController  = controller.NewLogController(logService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes
	routes.User(server, userController, jwtService)
	routes.Auth(server, userController, jwtService)
	routes.Log(server, logController, jwtService)
	routes.Management(server, userController, jwtService)

	server.Static("/assets", "./assets")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
