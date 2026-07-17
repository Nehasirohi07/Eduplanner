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

// CreateSubject godoc
// @Summary Create a subject
// @Description add a new subject to a course
// @Tags subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param subject body models.Subject true "Subject Details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id}/subjects [post]

func CreateSubject(w http.ResponseWriter, r *http.Request) {

	var subject models.Subject

	err := json.NewDecoder(r.Body).Decode(&subject)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	subject = utils.SanitizeSubject(subject)

	err = utils.ValidateSubject(subject)

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
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

	subject.CourseID = courseIDInt

	var courseExists int

	err = database.DB.QueryRow(
		"SELECT id FROM courses WHERE id =? AND user_id =?",
		courseIDInt,
		userID,
	).Scan(&courseExists)

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
		"INSERT INTO subjects(course_id, subject_name, daily_target_minutes) VALUES(? , ? , ?)",
		subject.CourseID,
		subject.SubjectName,
		subject.DailyTargetMinutes,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create subject",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"subject created successfully",
		nil,
	)
}
