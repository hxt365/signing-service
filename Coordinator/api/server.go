package api

import (
	"coordinator/usecase"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
	pu  *usecase.ProgressUseCase
}

func NewServer(pu *usecase.ProgressUseCase) *Server {
	s := &Server{
		mux: http.NewServeMux(),
		pu:  pu,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.Handle("/api/progresses", s.handleProgress())
}
