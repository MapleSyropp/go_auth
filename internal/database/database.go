package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/MapleSyropp/go_auth/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *Postgres
)

type Postgres struct {
	db *sql.DB
}

func CreateDatabase() (*Postgres, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var (
		database   = os.Getenv("DB_DATABASE")
		password   = os.Getenv("DB_PASSWORD")
		username   = os.Getenv("DB_USERNAME")
		dbInstance *Postgres
	)

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", username, database, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	dbInstance = &Postgres{
		db: db,
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	CreateUserTable(dbInstance)
	return dbInstance, nil
}

func CreateUserTable(postgres *Postgres) error {
	q := `CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		username text NOT NULL,
		password text NOT NULL
	);`
	_, err := postgres.db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func SaveUser(user *models.UserReq, postgres *Postgres) error {
	pstmt, err := postgres.db.Prepare("INSERT INTO Users (username, password) VALUES ($1, $2);")
	if err != nil {
		return err
	}
	defer pstmt.Close()

	_, err = pstmt.Exec(user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(username string, postgres *Postgres) (*models.User, error) {
	q := postgres.db.QueryRow("SELECT * FROM Users WHERE username=$1 LIMIT 1", username)
	newUser := new(models.User)
	err := q.Scan(&newUser.ID, &newUser.Username, &newUser.Password)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
