package models

import "time"

type Candidate struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Age         int       `json:"age" binding:"required"`
	Essay       string    `json:"essay" binding:"required"`
	Experience  string    `json:"experience" binding:"required"`
	Motivation  string    `json:"motivation" binding:"required"`
	Score       int       `json:"score"`
	Explanation string    `json:"explanation"`
	AIDetected  bool      `json:"ai_detected"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}