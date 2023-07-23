package service

import (
	"net/http"
	"testing"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMockNotesServiceLogin_SignUp(t *testing.T) {

	// Create a test UserSignUp model
	signUp := models.UserSignUp{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	login := models.UserLogin{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	mockDB := &db.MockPostgres{
		User:      &signUp,
		UserLogin: &login,
	}

	mockService := NewMockNotesService(mockDB)

	transactionID := uuid.New().String()

	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{transactionID},
			}},
	}

	err := mockService.SignUp(ctx, signUp)
	assert.NoError(t, err)

	err = mockService.Login(ctx, login)
	assert.NoError(t, err)

}

func TestMockNotesService_CreateGetDeleteNote(t *testing.T) {

	// Create a test Notes model
	note := models.Notes{
		Note:      "Test Note",
		SessionID: uuid.New().String(),
		NoteId:    "1",
	}

	mockDB := &db.MockPostgres{
		Note: &note,
	}

	mockService := NewMockNotesService(mockDB)

	transactionID := uuid.New().String()

	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{transactionID},
			}},
	}

	err := mockService.CreateNote(ctx, note)
	assert.NoError(t, err)

	err = mockService.GetNote(ctx, note)
	assert.NoError(t, err)

	err = mockService.DeleteNote(ctx, note)
	assert.NoError(t, err)
}
