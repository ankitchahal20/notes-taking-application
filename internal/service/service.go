package service

import (
	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/gin-gonic/gin"
)

var notesClient *NotesService

type NotesService struct {
	repo db.NotesDBService
}

func NewNotesService(conn db.NotesDBService) *NotesService {
	notesClient = &NotesService{
		repo: conn,
	}
	return notesClient
}

// This is a function to process the users details and subsequently using it for login.
func Login() func(ctx *gin.Context) {
	return func(context *gin.Context) {
	}
}

// This is a function to process the users details for sign-up and subsequently storing it in DB.
func SignUp() func(ctx *gin.Context) {
	return func(context *gin.Context) {
	}
}
