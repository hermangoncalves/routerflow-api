package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hermangoncalves/routerflow-api/features/auth"
	"github.com/hermangoncalves/routerflow-api/pkg/db"
)

func main() {
	db := db.InitDB()

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	auth.SetupRoutes(r, db)

	r.Run(":8080")
}
