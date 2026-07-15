package utils

import (
	"errors"
	"net/mail"
	"strings"

	"eduplanner/models"
)

func ValidateUser(user models.User) error {

	if strings.TrimSpace(user.Name) == "" {
		return errors.New("Name is requried")
	}

	if len(user.Name) < 3 {
		return errors.New("Name must be at least 3 characters")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("Email is required")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("Invalid Email Format ")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("Password is required")
	}

	if len(user.Password) < 6 {
		return errors.New("Password must be atleast 6 characters ")
	}

	return nil
}

func ValidateCourse(course models.Course) error {

	if strings.TrimSpace(course.CourseName) == "" {
		return errors.New("course name is required")
	}

	if len(course.CourseName) > 50 {
		return errors.New("course name must not exceed  50 characters")
	}

	return nil
}

func ValidateSubject(subject models.Subject) error {

	if strings.TrimSpace(subject.SubjectName) == "" {
		return errors.New("subject name is required")
	}

	if len(subject.SubjectName) > 30 {
		return errors.New("subject name must be in between 30 characters ")
	}

	return nil
}

func ValidateStudyGoal(studyGoal models.StudyGoal) error {

	if studyGoal.SubjectID <= 0 {
		return errors.New("Invalid subject")
	}

	if studyGoal.TargetMinutes <= 0 {
		return errors.New("Target minutes must be greater than 0")
	}

	if studyGoal.Deadline.IsZero() {
		return errors.New("Deadline is required")
	}

	if studyGoal.Status != "Pending" &&
		studyGoal.Status != "in_progress" &&
		studyGoal.Status != "complete" {

		return errors.New("Invalid status")
	}

	return nil
}
