package Dto

type LoginDto struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
