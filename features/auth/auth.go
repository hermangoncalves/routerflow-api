package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/login", loginHandler(db))
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var creds Credentials
		if err := ctx.ShouldBindJSON(&creds); err != nil {
			log.Printf("Invalid login request: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		var storedUser User
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", creds.Username).
			Scan(&storedUser.Username, &storedUser.Password)
		if err != nil || storedUser.Password != creds.Password {
			log.Println("Login failed for:", creds.Username)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		log.Println("Login successful for:", creds.Username)
		ctx.JSON(200, gin.H{"token": "fake-jwt-token"})
	}
}
