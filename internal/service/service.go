package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/ankit/project/notes-taking-application/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

var (
	notesClient *NotesService
)

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
	return func(ctx *gin.Context) {
		fmt.Println("Received request for user login !!!")
		var userLogin models.UserLogin
		if err := ctx.ShouldBindBodyWith(&userLogin, binding.JSON); err == nil {
			fmt.Println("User is trying to login !!!")
			_, err := notesClient.userLogin(ctx, userLogin)
			if err != nil {
				ctx.Writer.WriteHeader(err.Code)
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				return
			}

			sessionID := uuid.New().String()
			session, _ := utils.Store.Get(ctx.Request, sessionID)

			// Set session values
			session.Values["expiryTime"] = time.Now().Add(1 * time.Minute).Unix()

			// Save the session
			session.Values["sessionID"] = sessionID
			seesionErr := session.Save(ctx.Request, ctx.Writer)
			if seesionErr != nil {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": seesionErr.Error(),
				})
				return
			}

			fmt.Println("User Login is successful !!! ", session.Values["sessionID"])

			ctx.JSON(http.StatusOK, map[string]string{
				"sid": fmt.Sprintf("%v", session.Values["sessionID"]),
			})
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) userLogin(ctx *gin.Context, userLogin models.UserLogin) (string, *noteserror.NotesError) {
	if userLogin.Email == "" {
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is missing",
		}
	}
	if userLogin.Password == "" {
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "password is missing",
		}
	}

	pass, err := service.repo.Login(ctx, userLogin)
	if err != nil {
		return "", err
	}

	if pass == userLogin.Password {
		return "Login Successfull", nil
	}

	return "", &noteserror.NotesError{
		Code:    http.StatusBadRequest,
		Message: "incorrect password",
	}
}

// This is a function to process the users details for sign-up and subsequently storing it in DB.
func SignUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Received request for user signup !!!")
		var userSignUp models.UserSignUp
		if err := ctx.ShouldBindBodyWith(&userSignUp, binding.JSON); err == nil {
			fmt.Println("User is trying to signup !!!")
			err := notesClient.userSignUp(ctx, userSignUp)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}

			fmt.Println("User signup is successful !!!")
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) userSignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	if userSignUp.Email == "" {
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is missing",
		}
	}

	if userSignUp.Name == "" {
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "name is missing",
		}
	}

	if userSignUp.Password == "" {
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "password is missing",
		}
	}

	err := service.repo.SignUp(ctx, userSignUp)
	if err != nil {
		return err
	}

	return nil
}

func CreateNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Received request for user signup !!!")
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			fmt.Println("User is trying to signup !!!")
			notesId, err := notesClient.createNote(ctx, notes)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}
			fmt.Println("notesId : ", notesId)
			ctx.JSON(http.StatusOK, map[string]string{
				"id": notesId,
			})
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) createNote(ctx *gin.Context, notes models.Notes) (string, *noteserror.NotesError) {
	if notes.Note == "" {
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "note is empty",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	notesId, err := service.repo.CreateNotes(ctx, notes)
	if err != nil {
		return "", err
	}
	return notesId, nil
}

func DeleteNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Received request for a note deletion !!!")
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			fmt.Println("User is deelte the note to signup !!!")
			err := notesClient.deleteNote(ctx, notes)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}
			// ctx.JSON(http.StatusOK, map[string]string{
			// 	"id": fmt.Sprintf("%v", *notesId),
			// })
			fmt.Println("Deletion successfull")
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) deleteNote(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	if notes.NoteId == "" {
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "id is missing",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}

	err := service.repo.DeleteNotes(ctx, notes)
	if err != nil {
		return err
	}
	return nil
}

func GetNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println("Received request to get all notes !!!")
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			fmt.Println("User is trying to get all the note !!!")
			fetchedNotes, err := notesClient.getNotes(ctx, notes)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}
			ctx.JSON(http.StatusOK, map[string][]models.Notes{
				"notes": fetchedNotes,
			})

			fmt.Println("Get successfull")
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) getNotes(ctx *gin.Context, notes models.Notes) ([]models.Notes, *noteserror.NotesError) {

	fetchedNotes, err := service.repo.GetNotes(ctx)
	if err != nil {
		return []models.Notes{}, err
	}
	return fetchedNotes, nil
}
