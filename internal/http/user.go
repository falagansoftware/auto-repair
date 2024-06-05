package http

import (
	"log"
	"net/http"
	"strconv"

	autorepair "github.com/falagansoftware/auto-repair/internal"
	html "github.com/falagansoftware/auto-repair/internal/http/html"
)

// Routes

func (s *Server) registerUserRoutes() {
	s.router.HandleFunc("/users", s.handleUserList).Methods("GET")
	s.router.HandleFunc("/users/filter", s.handleUserListFilter).Methods("GET")
	s.router.HandleFunc("/users/search", s.handleUserListSearch).Methods("POST")
}

// Handlers

func (s *Server) handleUserList(w http.ResponseWriter, r *http.Request) {
	//Lang
	lang := r.URL.Query().Get("lang")
	// Filter
	filter := newUserFilters(r)
	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	log.Print(filter)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
	// Render Users
	view := html.UserList(users, filter.Sort, filter.Order, s.I18n.LangT(lang))
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

func (s *Server) handleUserListFilter(w http.ResponseWriter, r *http.Request) {
	//Lang
	lang := r.URL.Query().Get("lang")
	// Filter
	filter := newUserFilters(r)
	// Get all users
	users, err := s.UserService.FindUsers(r.Context(), filter)
	log.Print(filter)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
	// Render Users
	view := html.UserListSync(users, filter.Sort, filter.Order, s.I18n.LangT(lang))
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

func (s *Server) handleUserListSearch(w http.ResponseWriter, r *http.Request) {
	//Lang
	lang := r.URL.Query().Get("lang")
	r.ParseForm() // Parses the request body
	search := r.Form.Get("search")
	// Get all users
	users, err := s.UserService.FindUsersGlobally(r.Context(), &search)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
	// Render Users
	view := html.UserListSync(users, "", "", s.I18n.LangT(lang))
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
		return
	}
}

// Helpers

func newUserFilters(r *http.Request) *autorepair.UserFilters {
	// Filters
	filter := &autorepair.UserFilters{}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		filter.Limit = 20
	} else {
		filter.Limit = limit
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

	if err != nil {
		offset = 0
	} else {
		filter.Offset = offset
	}

	name := r.URL.Query().Get("name")

	if name != "" {
		filter.Name = &name
	}

	surname := r.URL.Query().Get("surname")

	if surname != "" {
		filter.Surname = &surname
	}

	email := r.URL.Query().Get("email")

	if email != "" {
		filter.Email = &email
	}

	active := r.URL.Query().Get("active")

	if active != "" {
		filter.Active, _ = strconv.ParseBool(active)
	}

	sort := r.URL.Query().Get("sort")

	if sort != "" {
		filter.Sort = sort
	}

	order := r.URL.Query().Get("order")

	if order != "" {
		filter.Order = order
	}

	return filter
}
