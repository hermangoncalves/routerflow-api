package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidShouldBindJSON(ctx *gin.Context, obj any) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		log.Printf("Invalid JSON request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return false
	}

	return true
}
