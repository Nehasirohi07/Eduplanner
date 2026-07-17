package handlers

import (
	"database/sql"
	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetSubjects godoc
// @Summary Get all subjects
// @Description Get all subjects of a course
// @Tags Subjects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses/{id}/subjects [get]

func GetSubjects(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
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

	_, err = database.DB.Query(
		"SELECT id FROM courses WHERE ID = ? AND user_id = ?",
		courseIDInt,
		userID,
	)

	if err == sql.ErrNoRows {
		utils.SendError(
			w,
			http.StatusForbidden,
			"Course not found or access denied",
		)
		return
	}

	var subjects []models.Subject

	rows, err := database.DB.Query(
		"SELECT id , course_id, subject_name, daily_target_minutes, created_at FROM subjects WHERE course_id = ?",
		courseIDInt,
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
		var subject models.Subject

		err = rows.Scan(
			&subject.ID,
			&subject.CourseID,
			&subject.SubjectName,
			&subject.DailyTargetMinutes,
			&subject.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				"Database error",
			)
			return
		}
		subjects = append(subjects, subject)
	}
	if err = rows.Err(); err != nil {
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
		"Subject fetched successfully",
		subjects,
	)

}
