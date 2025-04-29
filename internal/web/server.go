package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"task-dashboard/internal/service"
)

// Server manages the web server lifecycle and routes.
type Server struct {
	taskService *service.TaskService
	templates   *template.Template
	httpServer  *http.Server
	startTime   time.Time
}

// NewServer initializes a new Server instance with routes and templates.
func NewServer(taskService *service.TaskService, templatesDir string, port int) (*Server, error) {
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}

	mux := http.NewServeMux()

	s := &Server{
		taskService: taskService,
		templates:   templates,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/api/dashboard-data", s.handleDashboardData)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return s, nil
}

// Start runs the HTTP server. This call blocks until the server exits.
func (s *Server) Start() error {
	slog.Info("Server started", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown attempts a graceful shutdown of the HTTP server within the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// handleIndex renders the dashboard HTML page.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := s.taskService.GetDashboardData()
	if err != nil {
		http.Error(w, "Error fetching task data", http.StatusInternalServerError)
		return
	}
	if err := s.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		slog.Error("Template rendering failed", "error", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
	slog.Debug("Handling request", "method", r.Method, "url", r.URL.Path)
}

// handleDashboardData serves the dashboard data as JSON.
func (s *Server) handleDashboardData(w http.ResponseWriter, r *http.Request) {
	data, err := s.taskService.GetDashboardData()
	if err != nil {
		http.Error(w, "Error fetching task data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	slog.Debug("Handling request", "method", r.Method, "url", r.URL.Path)
}

// StartMetricsCollector launches a background ticker to collect metrics.
// Currently, this is just an expensive way to set an alarm that nobody hears.
func (s *Server) StartMetricsCollector(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				slog.Info("Metrics heartbeat", "uptime", time.Since(s.startTime))
			case <-ctx.Done():
				slog.Info("Stopping metrics collector")
				return
			}
		}
	}()
}
