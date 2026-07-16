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
		start_time
		)
		VALUES(?, NOW())`,
		subjectIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to start study Session",
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
	utils.SendSuccess(
		w,
		http.StatusOK,
		"Study session ended successfully",
		nil,
	)

}

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
		`SELECT subject.id
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
			"subject not found or access denied",
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

	var studySessions []models.StudySession

	rows, err := database.DB.Query(
		`SELECT
		id,
		subject_id,
		start_time,
		end_time,
		duration,
		created_at
		FROM study_session
		WHERE subject_id = ?`,
		subjectIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Database error",
		)
		return
	}

	defer rows.Close()

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
				"Database error",
			)
			return
		}

		studySessions = append(studySessions, studySession)

		if err = rows.Err(); err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Database error",
			)
			return
		}

	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Study sessions fetched successfully",
		studySessions,
	)

}
