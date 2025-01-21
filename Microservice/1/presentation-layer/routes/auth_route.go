package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mci-its/backend-service/domain-layer/middleware"
	"github.com/mci-its/backend-service/domain-layer/service"
	"github.com/mci-its/backend-service/presentation-layer/controller"
)

func Auth(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("/api/auth")
	{
		// User
		routes.POST("/register", userController.Register)
		routes.POST("/login", userController.Login)
		routes.DELETE("", middleware.Authenticate(jwtService), userController.Delete)
		routes.PATCH("/delete", middleware.Authenticate(jwtService), userController.SoftDelete)
		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
		routes.POST("/verify-email", userController.VerifyEmail)
		routes.POST("/send_verification_email", userController.SendVerificationEmail)
		routes.POST("/check-token", middleware.Authenticate(jwtService), userController.CheckToken)
		routes.POST("/refresh-token", middleware.Authenticate(jwtService), userController.RefreshToken)
		routes.PATCH("/change-password", middleware.Authenticate(jwtService), userController.ChangePassword)
		routes.POST("/send-reset-password", userController.SendResetPassword)
		routes.PATCH("/reset-password", userController.ResetPassword)
		routes.POST("/verify-otp", userController.VerifyOtp)
		routes.POST("/resend-otp", userController.ResendOtp)
		routes.PATCH("/upload-avatar", middleware.Authenticate(jwtService), userController.UploadAvatar)
	}
}
