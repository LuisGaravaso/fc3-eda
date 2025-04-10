package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort string
	Handlers      map[string]http.HandlerFunc
}

func NewWebServer(webServerPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: webServerPort,
		Handlers:      make(map[string]http.HandlerFunc),
	}
}

// For backward compatibility - defaults to POST
func (ws *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	ws.Handlers[path+":POST"] = handler
}

// Add specific method handlers
func (ws *WebServer) AddGetHandler(path string, handler http.HandlerFunc) {
	ws.Handlers[path+":GET"] = handler
}

func (ws *WebServer) AddPostHandler(path string, handler http.HandlerFunc) {
	ws.Handlers[path+":POST"] = handler
}

func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Logger)

	log.Println("Starting web server on port", ws.WebServerPort)

	for path, handler := range ws.Handlers {
		method := "POST" // Default method
		routePath := path

		// Check if the path has a method suffix
		if len(path) > 4 && path[len(path)-4] == ':' {
			method = path[len(path)-3:]
			routePath = path[:len(path)-4]
		}

		// Register the handler with the appropriate method
		switch method {
		case "GET":
			ws.Router.Get(routePath, handler)
		case "POST":
			ws.Router.Post(routePath, handler)
		default:
			ws.Router.Post(routePath, handler) // Fallback to POST
		}
	}

	http.ListenAndServe(ws.WebServerPort, ws.Router)
}
