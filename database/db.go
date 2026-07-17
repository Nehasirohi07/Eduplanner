package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(

		"%s:%s@tcp(%s:%s)/%s?parseTime=true",

		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	fmt.Println("User:", dbUser)
	fmt.Println("Host:", dbHost)
	fmt.Println("Port:", dbPort)
	fmt.Println("Database:", dbName)
	fmt.Println("Password Length:", len(dbPassword))
	fmt.Println("DSN:", dsn)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	for i := 1; i <= 10; i++ {

		err = DB.Ping()

		if err == nil {
			fmt.Println("✅ Database connected successfully!")
			return
		}
		fmt.Printf("Database not ready... Retry %d/10\n", i)

		time.Sleep(3 * time.Second)
	}

	log.Fatal("Could not connect to database after 10 retries")

}
func CreateTables() {

	userTable := `
	CREATE TABLE IF NOT EXISTS users(
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR (100) NOT NULL,
	email VARCHAR (100) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	role VARCHAR(20) DEFAULT 'user',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal("Error creating user tables:", err)
	}

	fmt.Println("✅ users table ready")

	coursesTable := `
CREATE TABLE IF NOT EXISTS courses(
id INT AUTO_INCREMENT PRIMARY KEY,
user_id INT NOT NULL,
FOREIGN KEY (user_id) REFERENCES users(id),
course_name VARCHAR(100)NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err = DB.Exec(coursesTable)
	if err != nil {
		log.Fatal("Error creating courses table:", err)
	}

	fmt.Println("✅ courses table ready")

	subjectsTable := `
CREATE TABLE IF NOT EXISTS subjects(
id INT AUTO_INCREMENT PRIMARY KEY,
course_id INT NOT NULL,
FOREIGN KEY (course_id) REFERENCES courses(id),
subject_name VARCHAR(100) NOT NULL,
daily_target_minutes INT NOT NULL, 
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err = DB.Exec(subjectsTable)
	if err != nil {
		log.Fatal("Error creating subject table:", err)
	}

	fmt.Println("✅ subject table ready")

	study_goalsTable := `
CREATE TABLE IF NOT EXISTS study_goals(
id INT AUTO_INCREMENT PRIMARY KEY,
subject_id INT NOT NULL,
FOREIGN KEY(subject_id) REFERENCES subjects(id),
target_minutes INT NOT NULL,
deadline DATE NOT NULL,
status VARCHAR(20) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err = DB.Exec(study_goalsTable)
	if err != nil {
		log.Fatal("Error creating study_goals table:", err)
	}

	fmt.Println("✅ study_goals ready")

	study_sessionsTable := `
CREATE TABLE IF NOT EXISTS study_sessions(
id INT AUTO_INCREMENT PRIMARY KEY,
subject_id INT NOT NULL,
FOREIGN KEY(subject_id) REFERENCES subjects(id),
start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
end_time TIMESTAMP DEFAULT NULL,
duration INT NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err = DB.Exec(study_sessionsTable)
	if err != nil {
		log.Fatal("Error creating study_session:", err)
	}

	fmt.Println("✅ study_session ready")

}
