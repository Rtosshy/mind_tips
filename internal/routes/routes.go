package routes

import (
	"database/sql"
	"net/http"

	"mind_tips/internal/controllers"
	"mind_tips/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, db *sql.DB) {
	authMiddleware := middleware.JWTMiddleware()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the blog app!"})
	})
	router.POST("/register", func(c *gin.Context) {
		controllers.RegisterUser(db, c)
	})
	router.POST("/login", func(c *gin.Context) {
		controllers.LoginUser(db, c)
	})
	router.GET("/user", authMiddleware, func(c *gin.Context) {
		controllers.GetUser(db, c)
	})
}
