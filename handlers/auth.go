package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"eduplanner/database"
	"eduplanner/models"
	"eduplanner/utils"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user = utils.SanitizeUser(user)

	err = utils.ValidateUser(user)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user.Password = string(hashedPassword)

	var existingID int

	err = database.DB.QueryRow(
		"SELECT id FROM users WHERE email = ?",
		user.Email,
	).Scan(&existingID)

	if err == nil {
		utils.SendError(w, http.StatusConflict, "Email already exists")
		return
	}

	if err != sql.ErrNoRows {
		utils.SendError(w, http.StatusInternalServerError, "Database error")
		return
	}

	_, err = database.DB.Exec(
		"INSERT INTO users(name,email,password)VALUES(? , ? , ?)",
		user.Name,
		user.Email,
		user.Password,
	)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	utils.SendSuccess(
		w,
		http.StatusCreated,
		"User registered successfully",
		nil,
	)

}

func Login(w http.ResponseWriter, r *http.Request) {

	var login models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User

	err = database.DB.QueryRow(
		"SELECT id, name , email, password, role FROM users WHERE email = ?",
		login.Email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		utils.SendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Database error")
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(login.Password),
	)

	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.SendSuccess(
		w,
		http.StatusOK,
		"Login successful",
		map[string]string{
			"token": token,
		},
	)
}
