package http

import (
	"log"
	"net/http"

	"github.com/falagansoftware/auto-repair/internal/http/html"
)

// Routes

func (s *Server) registerSignUpRoutes() {
	s.router.HandleFunc("/signup", s.handleSignUp).Methods("GET")
	s.router.HandleFunc("/signup", s.handleSignUpData).Methods("POST")
}

// Handlers

func (s *Server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	//Lang
	lang := r.URL.Query().Get("lang")
	// Render Users
	view := html.SignUp(s.i18n.LangT(lang))
	err := view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

func (s *Server) handleSignUpData(w http.ResponseWriter, r *http.Request) {
	// get post payload
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	// get form values
	name := r.FormValue("name")
	Surname := r.FormValue("surname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	// validate form values
	if name == "" || Surname == "" || email == "" || password == "" || confirmPassword == "" {
		log.Printf("Error: Empty form values")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}
	if password != confirmPassword {
		log.Printf("Error: Passwords do not match")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}
	// create user
	user := s.UserService.CreateUser(name, Surname, email, password)
	// redirect to login

	// Render Users
	view := html.SignUp(s.i18n.LangT(lang))
	err := view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}
