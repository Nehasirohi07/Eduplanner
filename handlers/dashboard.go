package handlers

import (
	"net/http"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	var dashboard models.Dashboard

	err := database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM courses 
		WHERE user_id = ?`,
		userID,
	).Scan(&dashboard.TotalCourses)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM subjects
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?`,
		userID,
	).Scan(&dashboard.TotalSubjects)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM study_goals
		JOIN subjects
		ON study_goals.subjects_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?
		AND study_goals.status = 'completed'`,
		userID,
	).Scan(&dashboard.CompletedGoals)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM study_goals
		JOIN subjects
		ON study_goals.subjects_id = subjects.id
		JOIN courses
		ON subjects.courses_id = courses.id
		WHERE courses.user_id = ?
		AND study_goals.status = 'pending'`,
		userID,
	).Scan(&dashboard.PendingGoals)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM study_sessions
		JOIN subjects
		ON study_sessions.subject_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?`,
		userID,
	).Scan(&dashboard.TotalStudyMinutes)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(duration),0)
		FROM study_sessions
		JOIN subjects
		ON study_session.subjects_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ? 
		AND DATE(start_time) = CURDATE()`,
		userID,
	).Scan(&dashboard.TodayStudyMinutes)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Dashboard fetched successfully",
		dashboard,
	)
}
