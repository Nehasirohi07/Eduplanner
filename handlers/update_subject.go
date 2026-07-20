package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"

	"github.com/gorilla/mux"
)

// UpdateSubject godoc
// @Summary Update a subject
// @Description update an existing subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Param subject body models.Subject true "Update Subject"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /subjects/{id} [put]

func UpdateSubject(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user ")
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

	var subject models.Subject

	err = json.NewDecoder(r.Body).Decode(&subject)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	subject = utils.SanitizeSubject(subject)

	err = utils.ValidateSubject(subject)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var existingSubject int

	err = database.DB.QueryRow(
		`SELECT subjects.id
		FROM subjects
		JOIN courses
		ON subjects.course_id = courses.id
		WHERE subjects.id = ?
		AND courses.user_id = ?`,
		subjectIDInt,
		userID,
	).Scan(&existingSubject)

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

	_, err = database.DB.Exec(
		`UPDATE subjects
		SET
			subject_name = ?,
			daily_target_minutes = ?
		WHERE id = ?`,
		subject.SubjectName,
		subject.DailyTargetMinutes,
		subjectIDInt,
	)
	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update subject",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Subject updated successfully",
		nil,
	)
}
