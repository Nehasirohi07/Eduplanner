package models

type Dashboard struct {
	TotalCourses      int `json:"total_courses"`
	TotalSubjects     int `json:"total_subjects"`
	CompletedGoals    int `json:"completed_goals"`
	PendingGoals      int `json:"pending_goals"`
	TotalStudyMinutes int `json:"total_study_minutes"`
	TodayStudyMinutes int `json:"today_study_minutes"`
}
