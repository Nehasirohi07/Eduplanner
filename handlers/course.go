package handlers

import (
	"encoding/json"
	"net/http"

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

}
