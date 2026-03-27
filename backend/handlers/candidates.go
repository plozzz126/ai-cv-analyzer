package handlers

import (
	"context"
	"net/http"
	"strconv"

	"invision-backend/db"
	"invision-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateCandidate(c *gin.Context) {
	var candidate models.Candidate

	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
	INSERT INTO candidates (name, age, essay, experience, motivation, score, explanation, ai_detected, status)
	VALUES ($1, $2, $3, $4, $5, 0, '', false, 'pending')
	RETURNING id, created_at
	`

	err := db.Conn.QueryRow(
		context.Background(),
		query,
		candidate.Name,
		candidate.Age,
		candidate.Essay,
		candidate.Experience,
		candidate.Motivation,
	).Scan(&candidate.ID, &candidate.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	candidate.Status = "pending"

	c.JSON(http.StatusOK, candidate)
}

func GetCandidates(c *gin.Context) {
	rows, err := db.Conn.Query(
		context.Background(),
		`SELECT id, name, age, essay, experience, motivation, score, explanation, ai_detected, status, created_at 
		 FROM candidates
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	candidates := make([]models.Candidate, 0)

	for rows.Next() {
		var candidate models.Candidate

		err := rows.Scan(
			&candidate.ID,
			&candidate.Name,
			&candidate.Age,
			&candidate.Essay,
			&candidate.Experience,
			&candidate.Motivation,
			&candidate.Score,
			&candidate.Explanation,
			&candidate.AIDetected,
			&candidate.Status,
			&candidate.CreatedAt,
		)

		if err != nil {
			println("SCAN ERROR:", err.Error())
			continue
		}

		candidates = append(candidates, candidate)
	}

	c.JSON(http.StatusOK, candidates)
}

func GetCandidateByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var candidate models.Candidate

	err = db.Conn.QueryRow(
		context.Background(),
		`SELECT id, name, age, essay, experience, motivation, score, explanation, ai_detected, status, created_at 
		 FROM candidates WHERE id=$1`,
		id,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.Age,
		&candidate.Essay,
		&candidate.Experience,
		&candidate.Motivation,
		&candidate.Score,
		&candidate.Explanation,
		&candidate.AIDetected,
		&candidate.Status,
		&candidate.CreatedAt,
	)

	if err != nil {
		c.JSON(404, gin.H{"error": "candidate not found"})
		return
	}

	c.JSON(200, candidate)
}

func ApproveCandidate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	_, err = db.Conn.Exec(
		context.Background(),
		"UPDATE candidates SET status='approved' WHERE id=$1",
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "approved"})
}

func RejectCandidate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	_, err = db.Conn.Exec(
		context.Background(),
		"UPDATE candidates SET status='rejected' WHERE id=$1",
		id,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "rejected"})
}

func GetLeaderboard(c *gin.Context) {
	rows, err := db.Conn.Query(
		context.Background(),
		`SELECT id, name, score, status FROM candidates 
		 ORDER BY score DESC LIMIT 10`,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type LeaderboardEntry struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Score  int    `json:"score"`
		Status string `json:"status"`
	}

	leaders := make([]LeaderboardEntry, 0)

	for rows.Next() {
		var l LeaderboardEntry
		err := rows.Scan(&l.ID, &l.Name, &l.Score, &l.Status)
		if err != nil {
			continue
		}
		leaders = append(leaders, l)
	}

	c.JSON(200, leaders)
}

func ScoreCandidate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	err = ScoreCandidateByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var candidate models.Candidate
	err = db.Conn.QueryRow(
		context.Background(),
		`SELECT id, name, age, essay, experience, motivation, score, explanation, ai_detected, status, created_at 
		 FROM candidates WHERE id=$1`,
		id,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.Age,
		&candidate.Essay,
		&candidate.Experience,
		&candidate.Motivation,
		&candidate.Score,
		&candidate.Explanation,
		&candidate.AIDetected,
		&candidate.Status,
		&candidate.CreatedAt,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "scored but failed to fetch result"})
		return
	}

	c.JSON(200, candidate)
}