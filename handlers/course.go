package handlers

import (
	"encoding/json"
	"net/http"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"
)

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a course for the logged-in user
// @Tags Courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course body models.Course true "Course details"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses [post]

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
		"INSERT INTO courses(user_id, course_name)VALUES(? , ?)",
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

// GetCourses godoc
// @Summary Get all courses
// @Description returns all courses of the logged-in user
// @Tags courses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /courses [get]

func GetCourses(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)

	if !ok {
		utils.SendError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	var courses []models.Course

	rows, err := database.DB.Query(
		"SELECT id, user_id, course_name, created_at FROM courses WHERE user_id = ?",
		userID,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var course models.Course

		err = rows.Scan(
			&course.ID,
			&course.UserID,
			&course.CourseName,
			&course.CreatedAt,
		)

		if err != nil {
			utils.SendError(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}

		courses = append(courses, course)

	}
	if err = rows.Err(); err != nil {
		utils.SendError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}
	utils.SendSuccess(
		w,
		http.StatusOK,
		"Courses fetched successfully",
		courses,
	)
}
