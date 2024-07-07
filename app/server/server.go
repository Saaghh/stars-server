package server

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"stars-server/app/config"

	"github.com/go-chi/chi/v5"
	"stars-server/app/generated/api-server"
	"stars-server/app/handlers"
)

type Server struct {
	APIServer *http.Server
}

func NewServer(cfg config.Config, proc handlers.ProcessorInterface) (*Server, error) {
	handler := handlers.NewHandler(proc)

	router := chi.NewRouter()

	server, err := api.NewServer(handler)
	if err != nil {
		zap.L().With(zap.Error(err)).Fatal("NewServer/api.NewServer(handler)")
	}

	router.Use(middleware.Logger)

	router.Mount("/api/", http.StripPrefix("/api", server))

	return &Server{
		APIServer: &http.Server{
			Addr:    cfg.ServerAddr,
			Handler: router,
		},
	}, nil
}

func (s *Server) runAPIServer() {
	zap.L().Info("Running API Server")
	if err := s.APIServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		zap.L().With(zap.Error(err)).Fatal("runAPIServer/ListenAndServe")
	}
}

func (s *Server) RunServer() {
	go s.runAPIServer()
}
