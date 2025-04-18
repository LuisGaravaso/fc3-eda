package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: port,
	}
}

func (ws *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	ws.Handlers[path] = handler
}

func (ws *WebServer) Start() error {
	ws.Router.Use(middleware.Logger)

	log.Println("Starting web server on port", ws.WebServerPort)

	for path, handler := range ws.Handlers {
		ws.Router.Post(path, handler)
	}

	return http.ListenAndServe(ws.WebServerPort, ws.Router)
}
