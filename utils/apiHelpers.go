package utils

import (
	"net/http"
	"encoding/json"
 	"github.com/hatim-lahwaouir/hypertube/types"
)

func WriteResp(w http.ResponseWriter,status int, v any) error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func MakeHandler(fc types.ApiFunc) http.HandlerFunc  {

	return func (w http.ResponseWriter, r *http.Request){
			if err := fc(w, r); err != nil{

				if e, ok := err.(types.ApiError); ok{
					WriteResp(w, e.Status, e)
				} else{
					errBody := types.ApiError{Err: "Internal Server Error", Status: http.StatusInternalServerError}
					WriteResp(w, errBody.Status, errBody)
				}

		}
	}
}
