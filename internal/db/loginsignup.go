package db

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/ankit/project/notes-taking-application/internal/utils"
	"github.com/gin-gonic/gin"
)

func (p postgres) SignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := `insert into users(name, password, emailid) values($1,$2,$3)`
	_, err := p.db.Exec(query, userSignUp.Name, userSignUp.Password, userSignUp.Email)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running insert query on users table, txid : %v", txid))
		if strings.Contains(err.Error(), "duplicate key value") {
			return &noteserror.NotesError{
				Trace:   txid,
				Code:    http.StatusBadRequest,
				Message: "product already added",
			}
		} else {
			return &noteserror.NotesError{
				Trace:   txid,
				Code:    http.StatusInternalServerError,
				Message: "unable to signup",
			}
		}
	}
	return nil
}

func (p postgres) Login(ctx *gin.Context, login models.UserLogin) (string, *noteserror.NotesError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := `SELECT password FROM users WHERE emailid=$1`
	var pass string
	if err := p.db.QueryRow(query, login.Email).Scan(&pass); err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running login db query, txid : %v", txid))
		return "", &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to get password",
			Trace:   txid,
		}
	}
	return pass, nil
}
