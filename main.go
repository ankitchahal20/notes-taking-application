package main

import (
	"log"

	"github.com/ankit/project/notes-taking-application/internal/config"
	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/server"
	"github.com/ankit/project/notes-taking-application/internal/service"
	"github.com/ankit/project/notes-taking-application/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Initializing the Log client
	utils.InitLogClient()

	utils.InitSession()

	// Initializing the GlobalConfig
	err := config.InitGlobalConfig()
	if err != nil {
		log.Fatalf("Unable to initialize global config")
	}

	// Establishing the connection to DB.
	postgres, err := db.New()
	if err != nil {
		log.Fatal("Unable to connect to DB : ", err)
	}

	// Initializing the client for notes service
	_ = service.NewNotesService(postgres)

	// Starting the server
	server.Start()
}
