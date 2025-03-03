package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hermangoncalves/routerflow-api/pkg/db"
)

func main() {
	_ = db.InitDB()

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
