package router

import (
	"net/http"

	"github.com/adrock-miles/about-me/internal/interfaces/http/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// New creates and configures the HTTP router.
func New(projectHandler *handler.ProjectHandler, contactHandler *handler.ContactHandler) http.Handler {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/projects", projectHandler.GetProjects).Methods("GET")
	api.HandleFunc("/contact", contactHandler.SendMessage).Methods("POST")

	// Serve static frontend files
	spa := spaHandler{staticPath: "web/dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	return c.Handler(r)
}

// spaHandler serves the React SPA, falling back to index.html for client-side routing.
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs := http.Dir(h.staticPath)

	// Try to serve the file directly
	if _, err := fs.Open(r.URL.Path); err != nil {
		// Fall back to index.html for SPA routing
		http.ServeFile(w, r, h.staticPath+"/"+h.indexPath)
		return
	}

	http.FileServer(fs).ServeHTTP(w, r)
}
