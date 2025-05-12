package db

import (
 	"github.com/hatim-lahwaouir/hypertube/types"
	"errors"
	"gorm.io/gorm"
)

func CreateUser(dto types.UserSignUpDto) error {
	user := User{
		Username:  dto.Username,
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password, // make sure the password is hashed 
	}
	result := DB.Create(&user)

	return  result.Error
}



func IsUserExists(username string , email string ) (bool , error) {
	
	result := DB.Where("username = ?", username).Or("email = ?",  email).First(&User{})

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil // user does not exist
	}

	if result.Error != nil {
		return false, result.Error // some other error occurred
	}

    	return true, nil 
}



func GetUser(username string , email string ) (*User, error) {
	
	var user User
	result := DB.Where("username = ?", username).Or("email = ?",  email).First(&user)

	if result.Error != nil {
		return nil, result.Error // some other error occurred
	}

    	return &user, nil 
}
