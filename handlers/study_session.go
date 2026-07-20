package handlers

import (
	"database/sql"
	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// StartStudySession godoc
// @Summary Start a new study session
// @Description Start a study session for a subject
// @Tags Study Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /subjects/{id}/study-session [post]

func StartStudySession(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	vars := mux.Vars(r)

	subjectID := vars["id"]

	subjectIDInt, err := strconv.Atoi(subjectID)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid subject ID",
		)
		return
	}

	var subjectExists int

	err = database.DB.QueryRow(
		`SELECT subjects.id
		FROM subjects
		JOIN courses
		ON subjects.course_id =courses.id
		WHERE subjects.id = ?
		AND courses.user_id = ?`,
		subjectIDInt,
		userID,
	).Scan(&subjectExists)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Subject not found or access denied",
		)
		return
	}

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	var activeSession int

	err = database.DB.QueryRow(
		`SELECT id
		FROM study_sessions
		WHERE subject_id = ?
		AND end_time IS NULL`,
		subjectIDInt,
	).Scan(&activeSession)

	if err == nil {
		utils.SendError(
			w,
			http.StatusConflict,
			"Study session already running",
		)
		return
	}

	if err != sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	_, err = database.DB.Exec(
		`INSERT INTO study_sessions(
		subject_id,
		start_time,
		duration
	)
	VALUES (?, NOW(), ?)`,
		subjectIDInt,
		0,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Study session started Successfully",
		nil,
	)

}

// EndStudySession godoc
// @Summary End study session
// @Description End the active study session of a subject
// @Tags Study Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /subjects/{id}/study-session [put]

func EndStudySession(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	vars := mux.Vars(r)

	subjectID := vars["id"]

	subjectIDInt, err := strconv.Atoi(subjectID)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid subject ID",
		)
		return
	}

	var subjectExists int

	err = database.DB.QueryRow(
		`SELECT subjects.id
		FROM subjects
		JOIN courses
		ON subjects.course_id =courses.id
		WHERE subjects.id = ?
		AND courses.user_id = ?`,
		subjectIDInt,
		userID,
	).Scan(&subjectExists)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Subject not found or access denied",
		)
		return
	}

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	var sessionID int
	var startTime time.Time

	err = database.DB.QueryRow(
		`SELECT id, start_time
		FROM study_sessions
		WHERE subject_id = ?
		AND end_time IS NULL`,
		subjectIDInt,
	).Scan(
		&sessionID,
		&startTime,
	)
	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"No active study session found",
		)
		return
	}

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}
	endTime := time.Now()

	duration := int(endTime.Sub(startTime).Minutes())

	_, err = database.DB.Exec(
		`UPDATE study_sessions
		SET end_time = ?, duration = ?
		WHERE id = ?`,
		endTime,
		duration,
		sessionID,
	)
	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to end study session",
		)
		return
	}

	var totalMinutes int

	err = database.DB.QueryRow(
		`SELECT COALESCE(SUM(duration), 0)
		FROM study_sessions
		WHERE subject_id = ?`,
		subjectIDInt,
	).Scan(&totalMinutes)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	var goalID int
	var targetMinutes int

	err = database.DB.QueryRow(
		`SELECT id, target_minutes
	FROM study_goals
	WHERE subject_id = ?
	ORDER BY created_at DESC
	LIMIT 1`,
		subjectIDInt,
	).Scan(&goalID, &targetMinutes)

	if err == sql.ErrNoRows {
		utils.SendSuccess(
			w,
			http.StatusOK,
			"Study session ended successfully",
			nil,
		)
		return
	}

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
		"Study session ended successfully",
		nil,
	)

}

// GetStudySession godoc
// @Summary Get study sessions
// @Description Get all study sessions of a subject
// @Tags Study Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /subjects/{id}/study-sessions [get]

func GetStudySession(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(
			w,
			http.StatusUnauthorized,
			"Invalid user",
		)
		return
	}

	vars := mux.Vars(r)

	subjectID := vars["id"]

	subjectIDInt, err := strconv.Atoi(subjectID)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid subject ID",
		)
		return
	}

	var subjectExists int

	err = database.DB.QueryRow(
		`SELECT subjects.id
		FROM subjects
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE subjects.id = ?
		AND courses.user_id = ?`,
		subjectIDInt,
		userID,
	).Scan(&subjectExists)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Subject not found or access denied",
		)
		return
	}

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(), // debugging
		)
		return
	}

	rows, err := database.DB.Query(
		`SELECT
			id,
			subject_id,
			start_time,
			end_time,
			duration,
			created_at
		FROM study_sessions
		WHERE subject_id = ?`,
		subjectIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(), // debugging
		)
		return
	}

	defer rows.Close()

	var studySessions []models.StudySession

	for rows.Next() {

		var studySession models.StudySession

		err = rows.Scan(
			&studySession.ID,
			&studySession.SubjectID,
			&studySession.StartTime,
			&studySession.EndTime,
			&studySession.Duration,
			&studySession.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				err.Error(), // debugging
			)
			return
		}

		studySessions = append(studySessions, studySession)
	}

	if err = rows.Err(); err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Study sessions fetched successfully",
		studySessions,
	)
}
