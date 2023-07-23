package db

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Set up the expected SQL query and result
	email := "test@example.com"
	password := "password123"
	mock.ExpectQuery(`SELECT password FROM notesusers WHERE emailid=\$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(password))

	// Create a new gin context for the test
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the Login function
	result, notesErr := p.Login(ctx, models.UserLogin{Email: email})
	if notesErr != nil {
		t.Fatalf("Login failed: %v", notesErr)
	}

	// Assert that the returned password matches the expected password
	if result != password {
		t.Errorf("Incorrect password. Expected %s, got %s", password, result)
	}

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestSignUp(t *testing.T) {
	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Set up the expected SQL query and result
	name := "John Doe"
	email := "test@example.com"
	password := "password123"
	mock.ExpectExec(`insert into notesusers\(name, password, emailid\) values\(\$1,\$2,\$3\)`).
		WithArgs(name, password, email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a new gin context for the test
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the SignUp function
	notesErr := p.SignUp(ctx, models.UserSignUp{Name: name, Email: email, Password: password})
	if notesErr != nil {
		t.Fatalf("SignUp failed: %v", notesErr)
	}

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
