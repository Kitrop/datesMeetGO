package models

import "time"

type UserLoginInput struct {
	Email 	 string `json:"email" validate:"required,email"` 	
	Password string `json:"password" validate:"required,min=6"`
}

type UserCreateInput struct {
	Username 	string  	`json:"username" validate:"required"`
	Email 	 	string		`json:"email" validate:"required,email"`
	Password 	string		`json:"password" validate:"required"`
	BirthDate time.Time `json:"birthDate" validate:"required"`
	Gender 		string		`json:"gender" validate:"required"` 
	Location 	string 		`json:"location" validate:"required"`
}