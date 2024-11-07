package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(db *sql.DB, c *gin.Context) {
	name, exists := c.Get("name")
	fmt.Println("Logged in as:", name)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authorized",
		})
		return
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var user User
	err := db.QueryRow("SELECT id, name FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User found", "user": user, "status": "success"})
}
