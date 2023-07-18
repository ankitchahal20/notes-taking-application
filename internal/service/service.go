package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankit/project/notes-taking-application/internal/db"
	"github.com/ankit/project/notes-taking-application/internal/models"
	noteserror "github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	store       sessions.CookieStore
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
			_, err := notesClient.UserLogin(ctx, userLogin)
			if err != nil {
				ctx.Writer.WriteHeader(err.Code)
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": err.Message,
				})
				return
			}
			store = sessions.NewCookieStore([]byte("hello"))
			session, _ := store.Get(ctx.Request, "session-name")

			// Set session values
			session.Values["userID"] = "user123"
			session.Values["expiryTime"] = time.Now().Add(60 * time.Second).Unix()

			// Save the session
			seesionErr := session.Save(ctx.Request, ctx.Writer)
			if seesionErr != nil {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				ctx.JSON(http.StatusOK, map[string]string{
					"Error": seesionErr.Error(),
				})
				return
			}
			fmt.Println("User Login is successful !!!")

			ctx.JSON(http.StatusOK, map[string]string{
				"sid": fmt.Sprintf("%v", session.Values["expiryTime"]),
			})
			ctx.Writer.WriteHeader(http.StatusOK)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) UserLogin(ctx *gin.Context, userLogin models.UserLogin) (string, *noteserror.NotesError) {
	fmt.Println("0", userLogin.Email, " : ", userLogin.Password)
	if userLogin.Email == "" {
		return "", &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is empty",
		}
	}
	fmt.Println("1")
	pass, err := service.repo.Login(ctx, userLogin)
	if err != nil {
		return "", err
	}
	fmt.Println("2")
	if pass == userLogin.Password {
		fmt.Println("3")
		return "Login Successfull", nil
	}
	fmt.Println("4")
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
			rowCount, err := notesClient.UserSignUp(ctx, userSignUp)
			if err != nil {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			if rowCount == 1 {
				fmt.Println("User signup is successful !!!")
				ctx.Writer.WriteHeader(http.StatusOK)
			} else {
				fmt.Println("User signup is unsuccessful !!!")
				ctx.Writer.WriteHeader(http.StatusUnauthorized)
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *NotesService) UserSignUp(ctx *gin.Context, userSignUp models.UserSignUp) (int, *noteserror.NotesError) {
	if userSignUp.Email == "" {
		return 0, &noteserror.NotesError{
			Code:    http.StatusBadRequest,
			Message: "EmailId is empty",
		}
	}

	err := service.repo.SignUp(ctx, userSignUp)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
