package db

import (
	"net/http"

	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MockNotesDBService interface {
	// login-signup
	Login(*gin.Context, models.UserLogin) (string, *noteserror.NotesError)
	SignUp(*gin.Context, models.UserSignUp) *noteserror.NotesError
	// notes
	CreateNotes(*gin.Context, models.Notes) (string, *noteserror.NotesError)
	DeleteNotes(*gin.Context, models.Notes) *noteserror.NotesError
	GetNotes(*gin.Context) ([]models.Notes, *noteserror.NotesError)
}

type MockPostgres struct {
	Note      *models.Notes
	User      *models.UserSignUp
	UserLogin *models.UserLogin
}

func (m *MockPostgres) Login(ctx *gin.Context, userLogin models.UserLogin) (string, *noteserror.NotesError) {
	if m.User.Password == userLogin.Password {
		return uuid.New().String(), nil
	}
	return "", &noteserror.NotesError{
		Code:    http.StatusBadGateway,
		Message: "password not found",
	}
}

func (m *MockPostgres) SignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	m.User.Email = userSignUp.Email
	m.User.Name = userSignUp.Name
	m.User.Password = userSignUp.Password
	return nil
}

func (m *MockPostgres) CreateNotes(ctx *gin.Context, notes models.Notes) (string, *noteserror.NotesError) {
	m.Note.Note = notes.Note
	noteId := "1"
	m.Note.NoteId = noteId
	return noteId, nil
}

func (m *MockPostgres) DeleteNotes(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	if m.Note.NoteId == notes.NoteId {
		m.Note = nil
	}
	return nil
}

func (m *MockPostgres) GetNotes(ctx *gin.Context) ([]models.Notes, *noteserror.NotesError) {
	fetchedNotes := []models.Notes{}
	fetchedNotes = append(fetchedNotes, *m.Note)
	return fetchedNotes, nil
}
