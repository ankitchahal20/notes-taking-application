package service

import (
	"fmt"

	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
)

var mockNotesClient *MockNotesService

type MockNotesService struct {
	MockRepo db.MockNotesDBService
	// Product  *models.Notes
	// User     *models.UserSignUp
}

func NewMockNotesService(conn db.MockNotesDBService) *MockNotesService {
	mockNotesClient = &MockNotesService{
		MockRepo: conn,
	}
	return mockNotesClient
}

func (m *MockNotesService) Login(ctx *gin.Context, login models.UserLogin) error {
	sesssionId, err := m.userLogin(ctx, login)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println("sesssionId : ", sesssionId)
	return nil
}

func (m *MockNotesService) userLogin(ctx *gin.Context, userLogin models.UserLogin) (string, *noteserror.NotesError) {
	sessionId, err := m.MockRepo.Login(ctx, userLogin)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func (m *MockNotesService) SignUp(ctx *gin.Context, userSignUp models.UserSignUp) error {
	err := m.userSignUp(ctx, userSignUp)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func (m *MockNotesService) userSignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	err := m.MockRepo.SignUp(ctx, userSignUp)
	if err != nil {
		return err
	}
	return nil
}

func (m *MockNotesService) CreateNote(ctx *gin.Context, note models.Notes) error {
	noteId, err := m.createNote(ctx, note)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println("NotesId : ", noteId)
	return nil
}

func (m *MockNotesService) createNote(ctx *gin.Context, notes models.Notes) (string, *noteserror.NotesError) {
	notesId, err := m.MockRepo.CreateNotes(ctx, notes)
	if err != nil {
		return "", err
	}

	return notesId, nil
}

func (m *MockNotesService) DeleteNote(ctx *gin.Context, note models.Notes) error {
	err := m.deleteNote(ctx, note)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func (m *MockNotesService) deleteNote(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	err := m.MockRepo.DeleteNotes(ctx, notes)
	return err
}

func (m *MockNotesService) GetNote(ctx *gin.Context, note models.Notes) error {
	notes, err := m.getNotes(ctx, note)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	fmt.Println("Notes : ", notes)
	return nil
}

func (m *MockNotesService) getNotes(ctx *gin.Context, notes models.Notes) ([]models.Notes, *noteserror.NotesError) {
	fetchedNotes, err := m.MockRepo.GetNotes(ctx)
	if err != nil {
		return nil, err
	}
	return fetchedNotes, nil
}
