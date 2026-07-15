package utils

import (
	"eduplanner/models"
	"strings"
)

func SanitizeUser(user models.User) models.User {

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Role = strings.TrimSpace(user.Role)
	user.Password = strings.TrimSpace(user.Password)

	return user

}

func SanitizeCourse(course models.Course) models.Course {

	course.CourseName = strings.TrimSpace(course.CourseName)

	return course
}

func SanitizeSubject(subject models.Subject) models.Subject {

	subject.SubjectName = strings.TrimSpace(subject.SubjectName)

	return subject
}

func SanitizeStudyGoal(studyGoal models.StudyGoal) models.StudyGoal {

	studyGoal.Status = strings.TrimSpace(strings.ToLower(studyGoal.Status))

	return studyGoal
}
