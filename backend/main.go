package main

import (
	"log"
	"os"

	"invision-backend/db"
	"invision-backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := db.InitDB()
	if err != nil {
		log.Fatal("DB error:", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/candidates", handlers.CreateCandidate)
	r.GET("/candidates", handlers.GetCandidates)
	r.GET("/candidates/:id", handlers.GetCandidateByID)

	r.POST("/candidates/:id/score", handlers.ScoreCandidate)

	r.POST("/candidates/:id/approve", handlers.ApproveCandidate)
	r.POST("/candidates/:id/reject", handlers.RejectCandidate)

	r.GET("/leaderboard", handlers.GetLeaderboard)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Сервер запущен на порту", port)
	r.Run(":" + port)
}