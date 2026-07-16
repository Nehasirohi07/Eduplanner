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

func CreateStudyGoal(w http.ResponseWriter, r *http.Request) {

	var studyGoal models.StudyGoal

	err := json.NewDecoder(r.Body).Decode(&studyGoal)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	studyGoal = utils.SanitizeStudyGoal(studyGoal)

	err = utils.ValidateStudyGoal(studyGoal)

	if err != nil {
		utils.SendError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
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
	}

	studyGoal.SubjectID = subjectIDInt

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
			"Subeject not found or access denied",
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
		`INSERT INTO study_goals(
			subject_id,
			target_minutes,
			deadline,
			status
        )
		VALUES(? , ? ,? ,?)`,
		studyGoal.SubjectID,
		studyGoal.TargetMinutes,
		studyGoal.Deadline,
		studyGoal.Status,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create study goal",
		)
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Study goal created successfully",
		nil,
	)

}
