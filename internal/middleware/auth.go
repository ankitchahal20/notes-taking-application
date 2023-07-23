package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/ankit/project/notes-taking-application/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// This function gets the unique transactionID
func getTransactionID(c *gin.Context) string {

	transactionID := c.GetHeader(constants.TransactionID)
	_, err := uuid.Parse(transactionID)
	if err != nil {
		transactionID = uuid.New().String()
		c.Request.Header.Set(constants.TransactionID, transactionID)
	}
	return transactionID
}

func AuthorizeUserRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// fetch the transactionID
		txid := getTransactionID(ctx)

		// validate the body params
		var notes models.Notes
		err := ctx.ShouldBindBodyWith(&notes, binding.JSON)
		if err != nil {
			utils.RespondWithError(ctx, http.StatusBadRequest, constants.InvalidBody)
			return
		}

		if notes.SessionID == "" {
			utils.RespondWithError(ctx, http.StatusBadRequest, "session id is missing")
			return
		}
		fmt.Println("notes : ", notes)

		notesError := authorizeUserRequest(ctx, txid, notes)
		if notesError != nil {
			utils.RespondWithError(ctx, notesError.Code, notesError.Message)
			return
		}
		ctx.Next()
	}
}

func authorizeUserRequest(ctx *gin.Context, txid string, user models.Notes) *noteserror.NotesError {
	// Get the session from the request
	session, _ := utils.Store.Get(ctx.Request, user.SessionID)
	// Check if the session ID is present
	if user.SessionID != session.Values["sessionID"] { // do this in loop
		// Session ID is invalid, user is not authenticated
		return &noteserror.NotesError{
			Code:    http.StatusUnauthorized,
			Message: "Session ID is invalid, user is not authenticated",
			Trace:   txid,
		}
	}

	// Check if the session has expired
	expiryTime, ok := session.Values["expiryTime"].(int64)
	if !ok || time.Now().Unix() > expiryTime {
		// Session has expired
		return &noteserror.NotesError{
			Code:    http.StatusUnauthorized,
			Message: "Session expired",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	return nil
}
