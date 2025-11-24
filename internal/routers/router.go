package routers

import (
	"go-web/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	engine := gin.Default()

	api := engine.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.RegisterHandler)
			auth.POST("/login", controller.LoginHandler)
			auth.POST("/logout", controller.LogoutHandler)
			auth.POST("/refresh", controller.RefreshTokenHandler)
			auth.POST("/forgot-password", controller.ForgotPasswordHandler)
			auth.POST("/reset-password", controller.ResetPasswordHandler)
			auth.POST("/check-username", controller.CheckUsernameHandler)
			auth.POST("/check-email", controller.CheckEmailHandler)
		}

		// 公开路由
		public := api.Group("/public")
		{
			public.GET("/health", controller.HealthHandler)
		}

		// 受保护的路由（需要认证）
		protected := api.Group("/protected")
		{
			users := protected.Group("/users")
			users.GET("/profile", controller.GetProfileHandler)
			users.PUT("/profile", controller.UpdateProfileHandler)
			users.PUT("/password", controller.ChangePasswordHandler)
			users.GET("", controller.GetUsersListHandler) // 获取用户列表（仅管理员）
		}
	}

	return engine
}
