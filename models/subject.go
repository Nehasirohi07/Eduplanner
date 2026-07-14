package models

import (
	"time"
)

type Subject struct {
	ID                 int       `json:"id"`
	CourseID           int       `json:"course_id"`
	SubjectName        string    `json:"subject_name"`
	DailyTargetMinutes int       `json:"daily_target_minutes"`
	CreatedAt          time.Time `json:"created_at"`
}
