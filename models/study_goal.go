package models

import (
	"time"
)

type StudyGoal struct {
	ID            int       `json:"id"`
	SubjectID     int       `json:"subject_id"`
	TargetMinutes int       `json:"target_minutes"`
	Deadline      time.Time `json:"deadline"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
