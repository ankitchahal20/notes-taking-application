package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/middleware"
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
		txid := middleware.GetTransactionID(ctx)
		utils.Logger.Info(fmt.Sprintf("received request for user login, txid : %v", txid))
		var userLogin models.UserLogin
		if err := ctx.ShouldBindBodyWith(&userLogin, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request is unmarshalled successfully, txid : %v", txid))
			_, err := notesClient.userLogin(ctx, userLogin)
			if err != nil {
				utils.Logger.Error(fmt.Sprintf("error received from service layer during user login, txid : %v", txid))
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
				utils.Logger.Error(fmt.Sprintf("error while saving the session, txid : %v", txid))
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": seesionErr.Error(),
				})
				return
			}

			utils.Logger.Info(fmt.Sprintf("user Login is successful, txid : %v", txid))

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
	txid := ctx.Request.Header.Get(constants.TransactionID)
	if userLogin.Email == "" {
		utils.Logger.Error(fmt.Sprintf("email id missing for user login, txid : %v", txid))
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is missing",
		}
	}
	if userLogin.Password == "" {
		utils.Logger.Error(fmt.Sprintf("password is missing for user login, txid : %v", txid))
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "password is missing",
		}
	}

	utils.Logger.Info(fmt.Sprintf("calling db layer for user login, txid : %v", txid))
	pass, err := service.repo.Login(ctx, userLogin)
	if err != nil {
		return "", err
	}

	if pass == userLogin.Password {
		return "Login Successfull", nil
	}

	utils.Logger.Error(fmt.Sprintf("incorrect password is used for user login, txid : %v", txid))
	return "", &noteserror.NotesError{
		Code:    http.StatusBadRequest,
		Message: "incorrect password",
	}
}

// This is a function to process the users details for sign-up and subsequently storing it in DB.
func SignUp() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := middleware.GetTransactionID(ctx)
		utils.Logger.Info(fmt.Sprintf("received request for user sign-up, txid : %v", txid))
		var userSignUp models.UserSignUp
		if err := ctx.ShouldBindBodyWith(&userSignUp, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request is unmarshalled successfully for user sign-up, txid : %v", txid))
			err := notesClient.userSignUp(ctx, userSignUp)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}
			utils.Logger.Info(fmt.Sprintf("user ser signup is successful, txid : %v", txid))
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) userSignUp(ctx *gin.Context, userSignUp models.UserSignUp) *noteserror.NotesError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	if userSignUp.Email == "" {
		utils.Logger.Error(fmt.Sprintf("email id missing for user sign-up, txid : %v", txid))
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is missing",
		}
	}

	if userSignUp.Name == "" {
		utils.Logger.Error(fmt.Sprintf("name is missing for user sign-up, txid : %v", txid))
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "name is missing",
		}
	}

	if userSignUp.Password == "" {
		utils.Logger.Error(fmt.Sprintf("password is missing for user sign-up, txid : %v", txid))
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "password is missing",
		}
	}

	utils.Logger.Info(fmt.Sprintf("calling db layer for user sign-up, txid : %v", txid))
	err := service.repo.SignUp(ctx, userSignUp)
	if err != nil {
		return err
	}

	return nil
}

func CreateNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for note creation, txid : %v", txid))
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for note creation is unmarshalled successfully, txid : %v", txid))
			notesId, err := notesClient.createNote(ctx, notes)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}
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
	txid := ctx.Request.Header.Get(constants.TransactionID)
	if notes.Note == "" {
		utils.Logger.Error(fmt.Sprintf("note is missing for user login, txid : %v", txid))
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "note is empty",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	utils.Logger.Info(fmt.Sprintf("calling db layer for note creation, txid : %v", txid))
	notesId, err := service.repo.CreateNotes(ctx, notes)
	if err != nil {
		return "", err
	}
	return notesId, nil
}

func DeleteNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for deleting a note, txid : %v", txid))
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request unmarshalled succesfully for deleting a note, txid : %v", txid))
			err := notesClient.deleteNote(ctx, notes)
			if err != nil {
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				ctx.Writer.WriteHeader(err.Code)
				return
			}

			utils.Logger.Info(fmt.Sprintf("user has successfully deleted a note, txid : %v", txid))
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) deleteNote(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	if notes.NoteId == "" {
		utils.Logger.Error(fmt.Sprintf("note id is missing for deleting a note, txid : %v", txid))
		return &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "id is missing",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	utils.Logger.Info(fmt.Sprintf("calling db layer for deleting a note, txid : %v", txid))
	err := service.repo.DeleteNotes(ctx, notes)
	if err != nil {
		return err
	}
	return nil
}

func GetNote() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for fetching all the notes, txid : %v", txid))
		var notes models.Notes
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request unmarshalled succesfully for fetching all the notes, txid : %v", txid))
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

			utils.Logger.Info(fmt.Sprintf("user has successfully fetched all the notes, txid : %v", txid))
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) getNotes(ctx *gin.Context, notes models.Notes) ([]models.Notes, *noteserror.NotesError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	utils.Logger.Info(fmt.Sprintf("calling db layer fetched all the notes, txid : %v", txid))
	fetchedNotes, err := service.repo.GetNotes(ctx)
	if err != nil {
		return []models.Notes{}, err
	}
	return fetchedNotes, nil
}
