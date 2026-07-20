package handlers

import (
	"database/sql"
	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// UpdateCourse godoc
// @Summary Update a course
// @Description Update an existing course
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param course body models.Course true "Update Course"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id} [put]

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
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

	var course models.Course

	err = json.NewDecoder(r.Body).Decode(&course)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	course = utils.SanitizeCourse(course)

	err = utils.ValidateCourse(course)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var existingCourse int

	err = database.DB.QueryRow(
		`SELECT id
		FROM courses
		WHERE id = ?
		AND user_id = ?`,
		courseIDInt,
		userID,
	).Scan(&existingCourse)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusForbidden,
			"course not found or access denied",
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
		`UPDATE courses
		SET course_name = ?
		WHERE id = ?`,
		course.CourseName,
		courseIDInt,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to update course",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Course update successfully",
		nil,
	)
}
