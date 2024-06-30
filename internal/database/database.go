package database

import (
	"database/sql"
	"fmt"
	"log"
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

func CreateDatabase() *Postgres {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err lodin", err)
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
		log.Fatal("failed to connect to db", err)
	}
	dbInstance = &Postgres{
		db: db,
	}
	if err := db.Ping(); err != nil {
		log.Fatal("no connection with db", err)
	}
	CreateUserTable(dbInstance)
	return dbInstance
}

func CreateUserTable(postgres *Postgres) {
	q := `CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		name text NOT NULL,
		password text NOT NULL
	);`
	_, err := postgres.db.Exec(q)
	if err != nil {
		log.Fatal("failed to create user table", err)
	}
}

func SaveUser(user *models.User, postgres *Postgres) error {
	pstmt, err := postgres.db.Prepare("INSERT INTO Users (name, password), VALUES ($1, $2);")
	defer pstmt.Close()

	u, err := pstmt.Exec(user.Name, user.Password)
	fmt.Println(u)
	return err
}
