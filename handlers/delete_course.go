package handlers

import (
	"database/sql"
	"eduplanner/database"
	"eduplanner/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Deletecourse godoc
// @Summary Delete a course
// @Description Delete an existing course
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id} [delete]

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
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

	courseID := vars["id"]

	courseIDInt, err := strconv.Atoi(courseID)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid course ID",
		)
		return
	}

	var existingCourse int

	err = database.DB.QueryRow(
		`SELECT id
		FROM courses
		WHERE id= ?
		AND user_id = ?`,
		courseIDInt,
		userID,
	).Scan(&existingCourse)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Course not found or access denied",
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
		`DELETE FROM courses
		WHERE id= ?`,
		courseIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to delete course",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Course deleted successfully",
		nil,
	)
}
