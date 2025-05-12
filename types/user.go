package types

type UserSignUpDto struct {
    	Username 	string		`json:"username" validate:"required,min=5,max=30"`
    	Email    	string		`json:"email" validate:"required,email"`
	FirstName	string		`json:"frist_name" validate:"required,min=3,max=30"`
	LastName	string		`json:"last_name" validate:"required,min=3,max=30"`
	Password	string		`json:"password" validate:"required,min=10,max=40"`
}


type UserSignInDto struct {
    	Email    	string 		`json:"email" validate:"required,email"`
    	Username 	string		`json:"username" validate:"required,min=5,max=30"`
	Password	string		`json:"password" validate:"required,min=10,max=40"`
}

