package types

import (
	"net/http"
)
type ApiMsg struct{
	Msg string	
	Status int
}

type ApiError struct{
	Err string	
	Status int
}



func (e ApiError) Error() string{
	return e.Err
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error
