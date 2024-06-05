package http

import (
	"encoding/json"
	"log"
	"net/http"

	autorepair "github.com/falagansoftware/auto-repair/internal"
	html "github.com/falagansoftware/auto-repair/internal/http/html"
)

// Routes

func (s *Server) registerNotificationsRoutes() {
	s.router.HandleFunc("/notifications", s.handleNotifications).Methods("POST")
}

// Handlers

func (s *Server) handleNotifications(w http.ResponseWriter, r *http.Request) {
	// get post payload
	var notification autorepair.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}
	// send notification
	view := html.Notification(notification.Level, notification.Message)
	err = view.Render(r.Context(), w)
	if err != nil {
		log.Printf("Internal Server Error: %v", err)
	}
}
