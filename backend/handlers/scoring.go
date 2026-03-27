package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"invision-backend/db"
	"invision-backend/models"
)

type AIRequest struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Essay      string `json:"essay"`
	Experience string `json:"experience"`
	Motivation string `json:"motivation"`
}

type AIResponse struct {
	Score       int    `json:"score"`
	Explanation string `json:"explanation"`
	AIDetected  bool   `json:"ai_detected"`
}

func ScoreCandidateByID(id int) error {
	var candidate models.Candidate

	err := db.Conn.QueryRow(
		context.Background(),
		"SELECT id, name, age, essay, experience, motivation FROM candidates WHERE id=$1",
		id,
	).Scan(
		&candidate.ID,
		&candidate.Name,
		&candidate.Age,
		&candidate.Essay,
		&candidate.Experience,
		&candidate.Motivation,
	)

	if err != nil {
		return fmt.Errorf("candidate not found: %w", err)
	}

	aiReq := AIRequest{
		ID:         candidate.ID,
		Name:       candidate.Name,
		Age:        candidate.Age,
		Essay:      candidate.Essay,
		Experience: candidate.Experience,
		Motivation: candidate.Motivation,
	}

	body, err := json.Marshal(aiReq)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		return fmt.Errorf("AI_SERVICE_URL не задан в .env")
	}

	resp, err := http.Post(aiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("AI сервис недоступен: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI сервис вернул статус %d", resp.StatusCode)
	}

	var aiResp AIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResp)
	if err != nil {
		return fmt.Errorf("ошибка парсинга ответа AI: %w", err)
	}

	_, err = db.Conn.Exec(
		context.Background(),
		`UPDATE candidates 
		 SET score=$1, explanation=$2, ai_detected=$3, status='scored'
		 WHERE id=$4`,
		aiResp.Score,
		aiResp.Explanation,
		aiResp.AIDetected,
		id,
	)

	return err
}