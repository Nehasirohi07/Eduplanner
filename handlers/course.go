package handlers

import (
	"encoding/json"
	"net/http"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
)

func CreateCourse(w http.ResponseWriter, r *http.Request) {

	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	course = utils.SanitizeCourse(course)

	err = utils.ValidateCourse(course)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	course.UserID = userID

	_, err = database.DB.Exec(
		"INSERT INTO course(user_id, course_name)VALUES(? , ?)",
		course.UserID,
		course.CourseName,
	)

	if err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			"Failed to create course",
		)
		return
	}
	utils.SendSuccess(
		w,
		http.StatusCreated,
		"Course created successfully",
		nil,
	)

}
