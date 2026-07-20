# 📚 EduPlanner

A RESTful backend API built with **Go (Golang)** for planning and tracking study progress. EduPlanner allows users to manage courses, subjects, study goals, study sessions, and view their study dashboard with secure JWT authentication.

---

## 🚀 Features

### 🔐 Authentication
- User Registration
- User Login
- JWT Authentication
- Protected Routes

### 📖 Course Management
- Create Course
- Get All Courses
- Update Course
- Delete Course

### 📘 Subject Management
- Create Subject
- Get Subjects
- Update Subject
- Delete Subject

### 🎯 Study Goals
- Create Study Goal
- Get Study Goals

### ⏱️ Study Sessions
- Start Study Session
- End Study Session
- Get Study Sessions
- Automatic Session Duration Calculation

### 📊 Dashboard
- View study statistics
- Track study progress

### ✅ Additional Features
- Input Validation
- Input Sanitization
- Authorization (User Ownership)
- Swagger API Documentation
- Docker Support
- MySQL Database

---

# 🛠 Tech Stack

- Go (Golang)
- Gorilla Mux
- MySQL
- JWT Authentication
- Swagger
- Docker
- REST API

---

# 📂 Project Structure

```
EduPlanner/
│
├── config/
├── database/
├── docs/
├── handlers/
├── middleware/
├── models/
├── routes/
├── utils/
├── .env
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

---

# ⚙️ Environment Variables

Create a `.env` file in the project root.

```env
DB_USER=root
DB_PASSWORD=root123
DB_HOST=mysql
DB_PORT=3306
DB_NAME=eduplanner

JWT_SECRET=your_secret_key

PORT=4040
```

---

# 🐳 Run with Docker

```bash
docker compose up --build
```

Backend:

```
http://localhost:4040
```

Swagger:

```
http://localhost:4040/swagger/index.html
```

---

# ▶️ Run Locally

### Install Dependencies

```bash
go mod tidy
```

### Generate Swagger

```bash
swag init
```

### Run Project

```bash
go run main.go
```

---

# 📄 API Documentation

Swagger UI

```
http://localhost:4040/swagger/index.html
```

---

# 📌 API Endpoints

## Authentication

| Method | Endpoint |
|---------|----------|
| POST | /register |
| POST | /login |

---

## Courses

| Method | Endpoint |
|---------|----------|
| POST | /courses |
| GET | /courses |
| PUT | /courses/{id} |
| DELETE | /courses/{id} |

---

## Subjects

| Method | Endpoint |
|---------|----------|
| POST | /courses/{id}/subjects |
| GET | /courses/{id}/subjects |
| PUT | /subjects/{id} |
| DELETE | /subjects/{id} |

---

## Study Goals

| Method | Endpoint |
|---------|----------|
| POST | /subjects/{id}/goals |
| GET | /subjects/{id}/goals |

---

## Study Sessions

| Method | Endpoint |
|---------|----------|
| POST | /subjects/{id}/study-session |
| PUT | /subjects/{id}/study-session |
| GET | /subjects/{id}/study-sessions |

---

## Dashboard

| Method | Endpoint |
|---------|----------|
| GET | /dashboard |

---

# 🔒 Authentication

All protected APIs require a JWT token.

```
Authorization: Bearer <your_jwt_token>
```

---

# 📦 Database

Main Tables

- users
- courses
- subjects
- study_goals
- study_sessions

---

# ✅ Project Highlights

- RESTful API Design
- JWT Authentication
- User Authorization
- Input Validation
- Clean Project Structure
- Dockerized Application
- Swagger Documentation
- MySQL Integration

---

# 🚀 Future Improvements

- Refresh Token Authentication
- Password Reset
- Email Verification
- Study Reminders
- Calendar Integration
- Analytics Dashboard
- Unit Tests
- CI/CD Pipeline

---

# 👩‍💻 Author

**Neha Sirohi**

B.Tech Computer Science Engineering

Backend Developer | Go Developer | AI & ML Enthusiast

GitHub:
https://github.com/Nehasirohi07

---

## ⭐ If you found this project useful, consider giving it a star!