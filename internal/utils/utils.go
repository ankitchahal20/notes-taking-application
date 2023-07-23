package utils

import (
	"github.com/ankit/project/notes-taking-application/internal/constants"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Store *sessions.CookieStore

func InitSession() {
	Store = sessions.NewCookieStore([]byte("hello"))
	Store.Options = &sessions.Options{
		MaxAge:   60, // Expiry time of 1 hour in seconds
		HttpOnly: true,
	}
}
func InitLogClient() {
	Logger, _ = zap.NewDevelopment()
}

func RespondWithError(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, noteserror.NotesError{
		Trace:   c.Request.Header.Get(constants.TransactionID),
		Code:    statusCode,
		Message: message,
	})
}
