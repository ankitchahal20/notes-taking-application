package models

type Note struct {
	ID   *int   `json:"id"`
	Note string `json:"note"`
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
