package main

import (
	"mind_tips/internal/database"
	"mind_tips/internal/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// MySQLに接続
	database.InitDB()

	router := gin.Default()
	db := database.GetDB()

	// ルーターのセットアップ
	routes.SetupRouter(router, db)

	router.Run(":8080")
}
