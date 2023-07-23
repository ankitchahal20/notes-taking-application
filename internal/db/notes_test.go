package db

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteNotes(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Set up the expected SQL query and result
	noteID := 1
	mock.ExpectExec(`DELETE FROM notes WHERE id=\$1`).
		WithArgs(fmt.Sprint(noteID)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new gin context for the test
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the DeleteNotes function
	notesErr := p.DeleteNotes(ctx, models.Notes{NoteId: fmt.Sprint(noteID)})
	if notesErr != nil {
		t.Fatalf("DeleteNotes failed: %v", notesErr)
	}

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestCreateNotes(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of the PostgreSQL repository
	p := postgres{
		db: db,
	}

	// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Set up the input notes
	notes := models.Notes{
		Note: "Test Note",
	}

	// Set up the expected SQL query and result
	mock.ExpectQuery(`insert into notes\(note\) values\(\$1\) RETURNING id`).
		WithArgs(notes.Note).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the function being tested
	notesID, notesErr := p.CreateNotes(ctx, notes)

	// Assert that the returned error is nil
	assert.Nil(t, notesErr)

	// Assert the expected notes ID
	expectedID := 1
	assert.Equal(t, expectedID, notesID)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetNotes(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of the PostgreSQL repository
	p := postgres{
		db: db,
	}

	// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Set up the expected SQL query and result
	mockRows := sqlmock.NewRows([]string{"id", "note"}).
		AddRow(1, "Note 1")

	mock.ExpectQuery(`SELECT id, note FROM notes`).
		WillReturnRows(mockRows)

	// Call the function being tested
	notes, notesErr := p.GetNotes(ctx)

	// Assert that the returned error is nil
	assert.Nil(t, notesErr)

	// Assert the expected number of notes
	expectedNotes := []models.Notes{
		{NoteId: "1", Note: "Note 1"},
	}
	assert.Equal(t, expectedNotes, notes)

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetNotes_Error(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of the PostgreSQL repository
	p := postgres{
		db: db,
	}

	// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Set up the expected SQL query to return an error
	mock.ExpectQuery(`SELECT id, note FROM notes`).
		WillReturnError(fmt.Errorf("database error"))

	// Call the function being tested
	notes, notesErr := p.GetNotes(ctx)

	// Assert that the returned error is not nil
	assert.NotNil(t, notesErr)

	// Assert the expected error message
	expectedError := &noteserror.NotesError{
		Code:    http.StatusInternalServerError,
		Message: "unable to get the notes",
		Trace:   "test-transaction-id",
	}
	assert.Equal(t, expectedError.Code, notesErr.Code)

	// Assert that the notes slice is nil
	assert.Nil(t, notes)

	// Assert that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
