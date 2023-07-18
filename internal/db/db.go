package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ankit/project/notes-taking-application/internal/config"
	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
)

type postgres struct{ db *sql.DB }

type NotesDBService interface {
	// login-signup
	Login(*gin.Context, models.UserLogin) (string, *noteserror.NotesError)
	SignUp(*gin.Context, models.UserSignUp) *noteserror.NotesError
	// notes
	CreateNotes(*gin.Context, models.Notes) (*int, *noteserror.NotesError)
	DeleteNotes(*gin.Context, models.Notes) *noteserror.NotesError
	GetNotes(*gin.Context, models.Notes) ([]models.Notes, *noteserror.NotesError)
}

func New() (postgres, error) {
	cfg := config.GetConfig()
	connString := "host=" + cfg.Database.Host + " " + "dbname=" + cfg.Database.DBname + " " + "password=" +
		cfg.Database.Password + " " + "user=" + cfg.Database.User + " " + "port=" + fmt.Sprint(cfg.Database.Port)

	conn, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to connect: %v\n", err))
		return postgres{}, err
	}

	log.Println("Connected to database")

	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot Ping the database")
		return postgres{}, err
	}
	log.Println("pinged database")

	return postgres{db: conn}, nil
}
