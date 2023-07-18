package models

type Notes struct {
	SessionID string `json:"sid,omitempty"`
	Note      string `json:"note,omitempty"`
	NoteId    string `json:"id,omitempty"`
}

type UserSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"` // FK
	Password string `json:"password"`
}
