package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	. "go_task4/handler"
	. "go_task4/model"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 全局错误处理中间件
	r.Use(errorHandler())
	r.Use(func(context *gin.Context) {
		context.Set("db", db)
	})
	api := r.Group("/api")
	api.POST("/register", Register)
	api.POST("/login", Login)
	{
		// 文章路由
		posts := api.Group("/posts")
		{
			posts.Use(authMiddleware()).GET("", ListPosts)
			posts.Use(authMiddleware()).GET("/:id", GetPost)
			posts.Use(authMiddleware()).POST("", CreatePost)
			posts.Use(authMiddleware()).PUT("/:id", UpdatePost)
			posts.Use(authMiddleware()).DELETE("/:id", DeletePost)
		}

		// 评论路由
		comments := api.Group("/posts/:id/comments")
		{
			comments.Use(authMiddleware()).GET("", ListComments)
			comments.Use(authMiddleware()).POST("", CreateComment)
		}
	}

	return r
}

// 认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			return
		}

		c.Set("userID", claims.UserID)
		fmt.Println("claims.UserID = ", claims.UserID)
		c.Next()
	}
}

// 错误处理中间件
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Err {
			case gorm.ErrRecordNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "资源不存在"})
			default:
				log.Printf("系统错误: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
			}
		}
	}
}
