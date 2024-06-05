package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	autorepair "github.com/falagansoftware/auto-repair/internal"
	"github.com/falagansoftware/auto-repair/internal/http/html"
	"github.com/falagansoftware/auto-repair/pkg/crypt"
)

// Routes

func (s *Server) registerSignUpRoutes() {
	s.router.HandleFunc("/signup", s.handleSignUp).Methods("GET")
	s.router.HandleFunc("/signup", s.handleSignUpData).Methods("POST")
}

// Handlers

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	// lang
	lang := r.URL.Query().Get("lang")
	// r ender Users
	view := html.SignUp(lang, s.I18n.LangT(lang))
	err := view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
	}
}

func (s *Server) handleSignUpData(w http.ResponseWriter, r *http.Request) {
	// lang
	lang := r.URL.Query().Get("lang")
	// get post payload
	var user autorepair.UserCreate
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	// validate
	err = s.Validator.Struct(user)
	if err != nil {
		log.Printf("Error validating user: %v", err)
	}
	// Hash password
	hash, err := crypt.HashPassword(user.Password)
	if err != nil {
		log.Printf("Error hashing user pass: %v", err)
	}
	user.Password = hash
	// create user
	err = s.UserService.CreateUser(r.Context(), &user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	// render response
	w.Header().Set("HX-Replace-Url", fmt.Sprintf("/signin?lang=%s", lang))
	view := html.SignUpSuccess(lang, s.I18n.LangT(lang))
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
	}

}
