package routes

import (
	"github.com/gin-gonic/gin"
	"golang_blog/controllers"
	"golang_blog/middleware"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()

	//使用中间件
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())

	//创建控制器实力
	authController := controllers.AuthController{}
	postController := controllers.PostController{}
	commentController := controllers.CommentController{}
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// 认证相关路由（无需认证）
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
	}
	//鉴权中间件
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// 用户信息
		authenticated.GET("/profile", authController.GetProfile)

		posts := authenticated.Group("/posts")
		{
			posts.POST("", postController.CreatePost)
			posts.PUT("/:id", postController.UpdatePost)
			posts.DELETE("/:id", postController.DeletePost)
		}
	}
	// 公开路由（无需认证）
	//public := api.Group("")
	{
		authenticated.GET("/posts", postController.GetPosts)
		authenticated.GET("/posts/:id", postController.GetPost)
	}
	comments := api.Group("/comments")
	{
		comments.GET("/post/:post_id", commentController.GetComments)
	}
	return r
}
