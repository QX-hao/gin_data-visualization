package router

import (
	"gin_data-visualization/api"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	// 创建Gin实例
	r := gin.Default()
	
	// 添加中间件
	r.Use(CORSMiddleware())
	r.Use(LoggerMiddleware())
	
	// 创建API路由组
	apiGroup := r.Group("/api/v1")
	{
		// 注册认证API
		authAPI := api.NewAuthAPI(db)
		authAPI.RegisterRoutes(apiGroup)
		
		// 健康检查路由
		apiGroup.GET("/health", authAPI.HealthCheck)
	}
	
	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	
	// 默认路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Gin Data Visualization API Server",
			"data": gin.H{
				"version": "1.0.0",
				"docs":    "/api/v1/docs/index.html",
			},
		})
	})
	
	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "接口不存在",
			"data":    nil,
		})
	})
	
	return r
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		// start := time.Now()
		
		// 处理请求
		c.Next()
		
		// 记录请求结束时间
		// latency := time.Since(start)
		
		// 这里可以添加日志记录逻辑
		// log.Printf("请求方法: %s, 路径: %s, 状态码: %d, 耗时: %v",
		// 	c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

// AuthMiddleware 认证中间件（预留）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在实际项目中，这里需要验证JWT token
		// token := c.GetHeader("Authorization")
		// if token == "" {
		// 	c.JSON(401, gin.H{"error": "未授权"})
		// 	c.Abort()
		// 	return
		// }
		
		// 验证token逻辑...
		
		c.Next()
	}
}