package routers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-routeros/routeros/v3"
	"github.com/hermangoncalves/routerflow-api/pkg/utils"
)

func SetupRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/routers/:id/status", statusHandler(db))
	r.POST("/routers/register", registerHandler(db))
}

func statusHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var router Router

		err := db.QueryRow("SELECT id, ip, username, password FROM routers WHERE id = ?", id).
			Scan(&router.ID, &router.IP, &router.Username, &router.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Router not found for ID:", id)
				ctx.JSON(404, gin.H{"error": "router not found"})
			} else {
				log.Println("Database error fetching router:", err)
				ctx.JSON(500, gin.H{"error": "internal server error"})
			}
			return
		}

		addr := fmt.Sprintf("%s:%s", router.IP, "8728")
		client, err := routeros.Dial(addr, router.Username, router.Password)
		if err != nil {
			log.Println("Failed to connect to router ID:", router.ID, "-", err)
			ctx.JSON(500, gin.H{"error": "connection failed"})
			return
		}
		defer client.Close()

		// Run command to get status
		reply, err := client.Run("/system/resource/print")
		if err != nil {
			log.Println("Failed to fetch status for router ID:", router.ID, "-", err)
			ctx.JSON(500, gin.H{"error": "command failed"})
			return
		}

		// Parse response into structured data
		status := RouterStatus{}
		for _, item := range reply.Re {
			for _, pair := range item.List {
				switch pair.Key {
				case "uptime":
					status.Uptime = pair.Value
				case "version":
					status.Version = pair.Value
				case "build-time":
					status.BuildTime = pair.Value
				case "factory-software":
					status.FactorySoftware = pair.Value
				case "free-memory":
					status.FreeMemory = pair.Value
				case "total-memory":
					status.TotalMemory = pair.Value
				case "cpu":
					status.CPU = pair.Value
				case "cpu-count":
					status.CPUCount = pair.Value
				case "cpu-frequency":
					status.CPUFrequency = pair.Value
				case "cpu-load":
					status.CPULoad = pair.Value
				case "free-hdd-space":
					status.FreeHDDSpace = pair.Value
				case "total-hdd-space":
					status.TotalHDDSpace = pair.Value
				case "write-sect-since-reboot":
					status.WriteSectSinceReboot = pair.Value
				case "write-sect-total":
					status.WriteSectTotal = pair.Value
				case "architecture-name":
					status.ArchitectureName = pair.Value
				case "board-name":
					status.BoardName = pair.Value
				case "platform":
					status.Platform = pair.Value
				}
			}
		}

		// Return structured JSON response
		ctx.JSON(http.StatusOK, status)
	}

}

func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req RouterRegisterRequest
		if !utils.ValidShouldBindJSON(ctx, &req) {
			return
		}

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM routers WHERE ip = ?)", req.IP).Scan(&exists)
		if err != nil {
			log.Println("Failed to check if router exists:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if exists {
			ctx.JSON(http.StatusConflict, gin.H{"error": "router with this IP already exists"})
			return
		}

		result, err := db.Exec("INSERT INTO routers (ip, username, password) VALUES (?, ?, ?)", req.IP, req.Username, req.Password)
		if err != nil {
			log.Println("Failed to register router:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Failed to get last inserted ID:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		log.Printf("Router registered with ID: %d\n", id)
		ctx.JSON(http.StatusCreated, gin.H{"message": "router registered", "id": id})

	}
}
