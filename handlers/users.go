package handlers



import (
	"net/http"
 	"github.com/hatim-lahwaouir/hypertube/utils"
	"fmt"
)

func Me(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value("user").(string)	
	fmt.Println(user)
	return utils.WriteResp(w,http.StatusOK, map[string]string{"me": "hello "})  

}

