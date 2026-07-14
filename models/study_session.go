package models

import (
	"time"
)

type StudySession struct {
	ID        int       `json:"id"`
	SubjectID int       `json:"subject_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}
