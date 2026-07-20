package handlers

import (
	"database/sql"
	"eduplanner/database"
	"eduplanner/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteSubject godoc
// @Summary Delete a subject
// @Description Delete an existing subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /Subjects/{id} [delete]

func DeleteSubject(w http.ResponseWriter, r *http.Request) {

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
	_, err = database.DB.Exec(
		`DELETE FROM subjects
		WHERE id = ?`,
		subjectIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete subject",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Subject deleted successfully",
		nil,
	)
}
