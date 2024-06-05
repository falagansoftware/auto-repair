package http

import (
	"encoding/json"
	"log"
	"net/http"

	autorepair "github.com/falagansoftware/auto-repair/internal"
	"github.com/falagansoftware/auto-repair/internal/http/html"
)

// Routes

func (s *Server) registerSignInRoutes() {
	s.router.HandleFunc("/signin", s.handleSignIn).Methods("GET")
	s.router.HandleFunc("/signin", s.handleSignInData).Methods("POST")
}

// Handlers

func (s *Server) handleSignIn(w http.ResponseWriter, r *http.Request) {
	//Lang
	lang := r.URL.Query().Get("lang")
	// Render Users
	view := html.SignIn(lang, s.I18n.LangT(lang))
	err := view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

func (s *Server) handleSignInData(w http.ResponseWriter, r *http.Request) {
	// get post payload
	var userLogin autorepair.UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	// validate
	err = s.Validator.Struct(userLogin)
	if err != nil {
		log.Printf("Error validating user: %v", err)
	}
	// get user
	// check hash password
	// hash, err := crypt.CheckPasswordHash(user.Password,)
	// if err != nil {
	// 	log.Printf("Error hashing user pass: %v", err)
	// }
	// user.Password = hash
	// // create user
	// err = s.UserService.CreateUser(r.Context(), &user)
	// if err != nil {
	// 	log.Printf("Error creating user: %v", err)
	// }
	// render Users
	// redirect to login
	// http.RedirectHandler("/signin", http.StatusSeeOther)
	// err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
	}
}
