package main

import (
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"time"
	"fmt"
	"github.com/go-playground/validator/v10"
	"errors"
)

// types

type apiFunc func(w http.ResponseWriter, r *http.Request) error


type apiMsg struct{
	Msg string	
	Status int
}

type apiError struct{
	Err string	
	Status int
}
// db functions

func Migrations(db *gorm.DB) {
	db.AutoMigrate(&User{})
	fmt.Println("✅ - migrations are done ! ")
}

func CreateDB() *gorm.DB{
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	} 

	fmt.Println("✅ - DB creation")
	return db
}


// User model
type User struct {
	ID		uint	`gorm:"primaryKey"` 
	Username	string
	Email		string
	FirstName	string
	LastName	string
	Password	string
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	gorm.DeletedAt `gorm:"index"` 
}

func (e apiError) Error() string{
	return e.Err
}

// Utils

func writeResp(w http.ResponseWriter,status int, v any) error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHandler(fc apiFunc) http.HandlerFunc  {

	return func (w http.ResponseWriter, r *http.Request){
			if err := fc(w, r); err != nil{

				if e, ok := err.(apiError); ok{
					writeResp(w, e.Status, e)
				} else{
					errBody := apiError{Err: "Internal Server Error", Status: http.StatusInternalServerError}
					writeResp(w, errBody.Status, errBody)
				}

		}
	}
}


var Validate *validator.Validate

func SetupValidator(){
	Validate = validator.New(validator.WithRequiredStructEnabled())
}	


// models 
type UserSignUpDto struct {
    	Username 	string		`json:"username" validate:"required,min=5,max=30"`
    	Email    	string		`json:"email" validate:"required,email"`
	FirstName	string		`json:"frist_name" validate:"required,min=3,max=30"`
	LastName	string		`json:"last_name" validate:"required,min=3,max=30"`
	Password	string		`json:"password" validate:"required,min=10,max=40"`
}


type UserSignInDto struct {
    	Email    	string 		`json:"email" validate:"required,email"`
	Password	string		`json:"password" validate:"required,min=10,max=40"`
}







func SignUp(w http.ResponseWriter, r *http.Request) error {
	// Validate.struct
	user := UserSignUpDto{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("❌ =>", err)
		return apiError{Err : "Invalid body" , Status : http.StatusBadRequest}
	}
	err = Validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			fmt.Println(err)
		}

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				fmt.Println(e.Field())
			}
		}

	}
	fmt.Println(user)

	return writeResp(w,201,apiMsg{Msg: "Hello ", Status:201})  
}


func SignIn(w http.ResponseWriter, r *http.Request) error {
	return writeResp(w,201,apiMsg{Msg: "Hello ", Status:201})  
}


func main() {
	SetupValidator()
	// creating of db and making migrations
	db := CreateDB()
	Migrations(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	r.Route("/api/auth",func(r chi.Router) {
 		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/sign-in", makeHandler(SignIn))
		r.Post("/sign-up", makeHandler(SignUp))
	})

	fmt.Println("✅ - Backend is running on port 3000")
	http.ListenAndServe(":3000", r)
}


