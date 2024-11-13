package routes

import (
	"database/sql"
	"net/http"

	"mind_tips/internal/controllers"
	"mind_tips/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, db *sql.DB) {
	authMiddleware := middleware.JWTMiddleware(db)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the blog app!"})
	})
	router.POST("/register", func(c *gin.Context) {
		controllers.Register(db, c)
	})
	router.POST("/login", authMiddleware.LoginHandler)

	userRoutes := router.Group("/user")
	userRoutes.Use(authMiddleware.MiddlewareFunc())
	{
		userRoutes.GET("", func(c *gin.Context) {
			controllers.GetUser(db, c)
		})
		userRoutes.DELETE("", func(c *gin.Context) {
			controllers.DeleteUser(db, c)
		})
		userRoutes.PUT("", func(c *gin.Context) {
			controllers.UpdateUser(db, c, authMiddleware)
		})
	}
}
