package routes

import (
	"database/sql"
	"net/http"

	"mind_tips/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, db *sql.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the blog app!"})
	})
	router.POST("/register", func(c *gin.Context) {
		controllers.RegisterUser(db, c)
	})
	router.POST("/login", func(c *gin.Context) {
		controllers.LoginUser(db, c)
	})
}
