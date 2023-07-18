package db

import (
	"net/http"
	"strings"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
)

func (p postgres) SignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	query := `insert into notesusers(name, password, emailid) values($1,$2,$3)`
	_, err := p.db.Exec(query, userSignUp.Name, userSignUp.Password, userSignUp.Email)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return &noteserror.NotesError{
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
				Code:    http.StatusBadRequest,
				Message: "product already added",
			}
		} else if strings.Contains(err.Error(), "violates foreign key constraint") {
			return &noteserror.NotesError{
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
				Code:    http.StatusBadRequest,
				Message: "user id is not found",
			}
		} else {
			return &noteserror.NotesError{
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
				Code:    http.StatusInternalServerError,
				Message: "unable to add product details",
			}
		}
	}
	return nil
}

func (p postgres) Login(ctx *gin.Context, login models.UserLogin) (string, *noteserror.NotesError) {
	query := `SELECT password FROM notesusers WHERE emailid=$1`
	var pass string
	if err := p.db.QueryRow(query, login.Email).Scan(&pass); err != nil {
		return "", &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to get password",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	return pass, nil
}
