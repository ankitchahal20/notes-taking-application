package db

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankit/project/notes-taking-application/internal/constants"
	"github.com/ankit/project/notes-taking-application/internal/models"
	"github.com/ankit/project/notes-taking-application/internal/noteserror"
	"github.com/ankit/project/notes-taking-application/internal/utils"
	"github.com/gin-gonic/gin"
)

func (p postgres) CreateNotes(ctx *gin.Context, notes models.Notes) (string, *noteserror.NotesError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := `insert into notes(note) values($1) RETURNING id`
	notesId := 0
	err := p.db.QueryRow(query, notes.Note).Scan(&notesId)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running insert query, txid : %v", txid))
		return "", &noteserror.NotesError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: "unable to add notes",
		}

	}
	id := strconv.Itoa(notesId)
	utils.Logger.Info(fmt.Sprintf("successfully added notes entry in db, txid : %v", txid))
	return id, nil
}

func (p postgres) DeleteNotes(ctx *gin.Context, notes models.Notes) *noteserror.NotesError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := `DELETE FROM notes WHERE id=$1`
	if _, err := p.db.Exec(query, notes.NoteId); err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running delete db query, txid : %v", txid))
		return &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to delete note",
			Trace:   txid,
		}
	}
	utils.Logger.Info(fmt.Sprintf("successfully deleted notes entry from db, txid : %v", txid))
	return nil
}

func (p postgres) GetNotes(ctx *gin.Context) ([]models.Notes, *noteserror.NotesError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	query := `SELECT id, note FROM notes`
	rows, err := p.db.Query(query)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("error while running select db query, txid : %v", txid))
		return nil, &noteserror.NotesError{
			Code:    http.StatusInternalServerError,
			Message: "unable to get the notes",
			Trace:   txid,
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
			utils.Logger.Error(fmt.Sprintf("error while scanning notes from db, txid : %v", txid))
			return nil, &noteserror.NotesError{
				Code:    http.StatusInternalServerError,
				Message: "unable to get the notes",
				Trace:   txid,
			}
		}
		scannedNotes = append(scannedNotes, note)
	}
	utils.Logger.Info(fmt.Sprintf("successfully fetched all the notes from db, txid : %v", txid))
	return scannedNotes, nil
}
