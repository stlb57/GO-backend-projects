package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Shortcut struct {
	ID        string    `json:"id"`
	Key       int       `json:"key"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Target    string    `json:"target"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateShortcutRequest struct {
	Key    int    `json:"key"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Target string `json:"target"`
}

type UpdateShortcutRequest struct {
	Name   *string `json:"name,omitempty"`
	Type   *string `json:"type,omitempty"`
	Target *string `json:"target,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Store struct {
	mu        sync.RWMutex
	shortcuts map[string]Shortcut
	nextID    int
}

func NewStore() *Store {
	return &Store{
		shortcuts: make(map[string]Shortcut),
		nextID:    1,
	}
}

func (s *Store) Create(shortcut Shortcut) Shortcut {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortcut.ID = strconv.Itoa(s.nextID)
	s.nextID++

	s.shortcuts[shortcut.ID] = shortcut

	return shortcut
}

func (s *Store) Get(id string) (Shortcut, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	shortcut, ok := s.shortcuts[id]

	return shortcut, ok
}

func (s *Store) List() []Shortcut {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Shortcut, 0, len(s.shortcuts))

	for _, shortcut := range s.shortcuts {
		result = append(result, shortcut)
	}

	return result
}

func (s *Store) Update(id string, shortcut Shortcut) (Shortcut, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.shortcuts[id]

	if !ok {
		return Shortcut{}, false
	}

	s.shortcuts[id] = shortcut

	return shortcut, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.shortcuts[id]

	if !ok {
		return false
	}

	delete(s.shortcuts, id)

	return true
}

func validateShortcut(key int, name, shortcutType, target string) error {
	var errs []string

	if key <= 0 {
		errs = append(errs, "key must be positive")
	}

	if strings.TrimSpace(name) == "" {
		errs = append(errs, "name cannot be empty")
	}

	switch shortcutType {
	case "url", "folder", "app":
	default:
		errs = append(errs, "type must be url, folder, or app")
	}

	if strings.TrimSpace(target) == "" {
		errs = append(errs, "target cannot be empty")
	}

	if shortcutType == "url" {
		parsed, err := url.ParseRequestURI(target)

		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			errs = append(errs, "target must be a valid URL")
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

type Server struct {
	store *Store
}

func (s *Server) shortcutsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listShortcuts(w, r)
	case http.MethodPost:
		s.createShortcut(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
	}
}

func (s *Server) shortcutHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(
		strings.Trim(r.URL.Path, "/"),
		"/",
	)

	if len(parts) != 4 {
		writeError(
			w,
			http.StatusNotFound,
			"NOT_FOUND",
			"shortcut not found",
		)
		return
	}

	id := parts[3]

	switch r.Method {
	case http.MethodGet:
		s.getShortcut(w, r, id)
	case http.MethodPatch:
		s.updateShortcut(w, r, id)
	case http.MethodDelete:
		s.deleteShortcut(w, r, id)
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeError(
			w,
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"method not allowed",
		)
	}
}

func (s *Server) createShortcut(w http.ResponseWriter, r *http.Request) {
	var req CreateShortcutRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid JSON body",
		)
		return
	}

	err = validateShortcut(
		req.Key,
		req.Name,
		req.Type,
		req.Target,
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			err.Error(),
		)
		return
	}

	shortcut := Shortcut{
		Key:       req.Key,
		Name:      req.Name,
		Type:      req.Type,
		Target:    req.Target,
		CreatedAt: time.Now(),
	}

	shortcut = s.store.Create(shortcut)

	writeJSON(
		w,
		http.StatusCreated,
		shortcut,
	)
}

func (s *Server) listShortcuts(w http.ResponseWriter, r *http.Request) {
	shortcuts := s.store.List()

	filterType := r.URL.Query().Get("type")

	if filterType != "" {
		filtered := make([]Shortcut, 0)

		for _, shortcut := range shortcuts {
			if shortcut.Type == filterType {
				filtered = append(filtered, shortcut)
			}
		}

		shortcuts = filtered
	}

	writeJSON(
		w,
		http.StatusOK,
		shortcuts,
	)
}

func (s *Server) getShortcut(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	shortcut, ok := s.store.Get(id)

	if !ok {
		writeError(
			w,
			http.StatusNotFound,
			"SHORTCUT_NOT_FOUND",
			"shortcut not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		shortcut,
	)
}

func (s *Server) updateShortcut(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	shortcut, ok := s.store.Get(id)

	if !ok {
		writeError(
			w,
			http.StatusNotFound,
			"SHORTCUT_NOT_FOUND",
			"shortcut not found",
		)
		return
	}

	var req UpdateShortcutRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid JSON body",
		)
		return
	}

	if req.Name != nil {
		shortcut.Name = *req.Name
	}

	if req.Type != nil {
		shortcut.Type = *req.Type
	}

	if req.Target != nil {
		shortcut.Target = *req.Target
	}

	err = validateShortcut(
		shortcut.Key,
		shortcut.Name,
		shortcut.Type,
		shortcut.Target,
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			err.Error(),
		)
		return
	}

	updated, _ := s.store.Update(id, shortcut)

	writeJSON(
		w,
		http.StatusOK,
		updated,
	)
}

func (s *Server) deleteShortcut(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	deleted := s.store.Delete(id)

	if !deleted {
		writeError(
			w,
			http.StatusNotFound,
			"SHORTCUT_NOT_FOUND",
			"shortcut not found",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	store := NewStore()

	server := &Server{
		store: store,
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/v1/shortcuts",
		server.shortcutsHandler,
	)

	mux.HandleFunc(
		"/api/v1/shortcuts/",
		server.shortcutHandler,
	)

	fmt.Println("Shortcut API running on :8080")

	err := http.ListenAndServe(
		":8080",
		mux,
	)

	if err != nil {
		fmt.Println(err)
	}
}
