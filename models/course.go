package models

import (
	"time"
)

type Course struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CourseName string    `json:"course_name"`
	CreateAt   time.Time `json:"created_at"`
}
