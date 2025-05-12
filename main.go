package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"fmt"
 	"github.com/hatim-lahwaouir/hypertube/handlers"
 	"github.com/hatim-lahwaouir/hypertube/utils"
 	"github.com/hatim-lahwaouir/hypertube/db"
	"context"
)


func AuthMiddleware(next http.Handler) http.Handler {
  	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len (tokenString) <= len("Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
		    	return
		}
		tokenString = tokenString[len("Bearer "):]

		id, ok := utils.VerifyToken(tokenString)

		if ok != true{
			w.WriteHeader(http.StatusUnauthorized)
		    	return
		}

		ctx := context.WithValue(r.Context(), "user", id)
		next.ServeHTTP(w, r.WithContext(ctx))
  })
}


func main() {
	// creating of db and making migrations

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	db.New()
	r.Route("/api/auth",func(r chi.Router) {
 		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/sign-in", utils.MakeHandler(handlers.SignIn))
		r.Post("/sign-up", utils.MakeHandler(handlers.SignUp))
	})

	r.Route("/api/user",func(r chi.Router) {
 		r.Use(AuthMiddleware)
		r.Get("/me", utils.MakeHandler(handlers.Me))
	})


	fmt.Println("✅ - Backend is running on port 3000")
	http.ListenAndServe(":3000", r)
}


