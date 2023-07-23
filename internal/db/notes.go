package db

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/gin-gonic/gin"
)

func (p postgres) CreateNotes(ctx *gin.Context, notes models.Notes) (string, *noteserror.NotesError) {
	fmt.Println("Inside UserSignUp : 0", notes)

	query := `insert into notes(note) values($1) RETURNING id`
	fmt.Println("Query : ", query, notes.Note)
	notesId := 0
	err := p.db.QueryRow(query, notes.Note).Scan(&notesId)
	if err != nil {

		return "", &noteserror.NotesError{
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
			Code:    http.StatusInternalServerError,
			Message: "unable to add notes",
		}

	}
	id := strconv.Itoa(notesId)
	return id, nil
}

func (p postgres) DeleteNotes(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	query := `DELETE FROM notes WHERE id=$1`

	if _, err := p.db.Exec(query, notes.NoteId); err != nil {
		return &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to delete note",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}

	return nil
}

func (p postgres) GetNotes(ctx *gin.Context) ([]models.Notes, *noteserror.NotesError) {
	query := `SELECT id, note FROM notes`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the notes",
			Trace:   ctx.Request.Header.Get(constants.TransactionID),
		}
	}
	defer rows.Close()

	// Create a slice to store the retrieved notes
	scannedNotes := []models.Notes{}

	// Iterate over the rows and scan the values into Note structs
	for rows.Next() {
		var note models.Notes
		err := rows.Scan(&note.NoteId, &note.Note)
		if err != nil {
			return nil, &noteserror.NotesError{
				Code:    http.StatusInternalServerError,
				Message: "unable to get the notes",
				Trace:   ctx.Request.Header.Get(constants.TransactionID),
			}
		}
		scannedNotes = append(scannedNotes, note)
	}
	return scannedNotes, nil
}
