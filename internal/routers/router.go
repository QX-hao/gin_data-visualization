package routers

import (
	"go-web/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	engine := gin.Default()

	// 添加CORS中间件
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	api := engine.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.RegisterHandler)
			auth.POST("/login", controller.LoginHandler)
			// auth.POST("/login", func(c *gin.Context) {
			// 	// 直接返回一个 401 错误，不进行任何数据库查询或逻辑判断
			// 	c.JSON(401, gin.H{
			// 		"code":    401,
			// 		"message": "这是一个直接返回的401测试",
			// 	})
			// })

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
