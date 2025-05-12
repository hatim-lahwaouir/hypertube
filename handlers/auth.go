package handlers



import (
	"net/http"
	"encoding/json"
	"fmt"
 	"github.com/hatim-lahwaouir/hypertube/types"
 	"github.com/hatim-lahwaouir/hypertube/utils"
 	"github.com/hatim-lahwaouir/hypertube/db"
)





func SignUp(w http.ResponseWriter, r *http.Request) error {
	dto := types.UserSignUpDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Invalid body" , Status : http.StatusBadRequest}
	}
	err = utils.Validate.Struct(dto)
	if err != nil {
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "data doesn't respect the rules body" , Status : http.StatusBadRequest}
	}


	var userExists bool
	userExists, err = db.IsUserExists(dto.Username, dto.Email)
	
	if err != nil{
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Server Error" , Status : http.StatusInternalServerError}
	}

	if userExists{
		return types.ApiError{Err : "Username or Email already exists" , Status : http.StatusBadRequest}
	}


	// password Hashing
	dto.Password, err  = utils.HashPassword(dto.Password)

	if err != nil{
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Server Error" , Status : http.StatusInternalServerError}
	}

	// creating the user
	if err = db.CreateUser(dto); err != nil {
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Server Error" , Status : http.StatusInternalServerError}
	}

	return utils.WriteResp(w,201,types.ApiMsg{Msg: "Accont created ", Status: http.StatusCreated})  
}


func SignIn(w http.ResponseWriter, r *http.Request) error {
	dto := types.UserSignInDto{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Invalid body" , Status : http.StatusBadRequest}
	}
	err = utils.Validate.Struct(dto)
	if err != nil {
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "data doesn't respect the rules body" , Status : http.StatusBadRequest}
	}


	var user *db.User
	user, err = db.GetUser(dto.Username, dto.Email)

	if err != nil{
		return types.ApiError{Err : "Invalid Request" , Status : http.StatusBadRequest}
	}
	// func VerifyPassword(password, hash string) bool {
	if utils.VerifyPassword(dto.Password, user.Password) == false{
		return types.ApiError{Err : "Invalid Passowrd or Email or Username" , Status : http.StatusBadRequest}
	}

	// func CreateToken(id uint64) (string, error) {
	var token string
	token, err  = utils.CreateToken(user.ID)
	if err != nil{
		fmt.Println("❌ =>", err)
		return types.ApiError{Err : "Server Error" , Status : http.StatusInternalServerError}
	}

	return utils.WriteResp(w,http.StatusOK, map[string]string{"access_token": token})  
}
