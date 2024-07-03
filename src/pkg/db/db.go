package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	config       RepositoryConfig
	dbConnection *sqlx.DB
}

type RepositoryConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewRepository(config RepositoryConfig) *Repository {
	return &Repository{
		config: config,
	}
}

func (r *Repository) ConnectToDB() {
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", r.config.User, r.config.DBName)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Connecting to DB failed", err)
	}

	r.dbConnection = db
}

// functions to get template by name etc
//tx := r.dbConnection.MustBegin()
// tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")

func (r *Repository) GetTemplateByName(name string) string {
	return ""
}
