package handlers

import (
	"net/http"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
)

// GetDashboard godoc
// @Summary Get dashboard
// @Description Get dashboard statistics for the logged-in user
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /dashboard [get]

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

	// Total Courses
	err := database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM courses
		WHERE user_id = ?`,
		userID,
	).Scan(&dashboard.TotalCourses)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Total Subjects
	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM subjects
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?`,
		userID,
	).Scan(&dashboard.TotalSubjects)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Completed Goals
	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM study_goals
		JOIN subjects
		ON study_goals.subject_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?
		AND study_goals.status = 'completed'`,
		userID,
	).Scan(&dashboard.CompletedGoals)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Pending Goals
	err = database.DB.QueryRow(
		`SELECT COUNT(*)
		FROM study_goals
		JOIN subjects
		ON study_goals.subject_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?
		AND study_goals.status = 'pending'`,
		userID,
	).Scan(&dashboard.PendingGoals)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Total Study Minutes
	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(duration),0)
		FROM study_sessions
		JOIN subjects
		ON study_sessions.subject_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?`,
		userID,
	).Scan(&dashboard.TotalStudyMinutes)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Today's Study Minutes
	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(duration),0)
		FROM study_sessions
		JOIN subjects
		ON study_sessions.subject_id = subjects.id
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE courses.user_id = ?
		AND DATE(start_time) = CURDATE()`,
		userID,
	).Scan(&dashboard.TodayStudyMinutes)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Dashboard fetched successfully",
		dashboard,
	)
}
