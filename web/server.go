package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/korchasa/dostavki/locator"
	"net/http"
)

type Server struct {
	router *mux.Router
	zakazClient *locator.ZakazClient
}

func New() *Server {
	cl := &http.Client{}
	zk := locator.New(cl)
	return &Server{
		router: mux.NewRouter(),
		zakazClient: zk,
	}
}

func (s *Server) Start(addr string) error {
	s.routes()
	if err := http.ListenAndServe(addr, s.router); err != nil {
		return fmt.Errorf("can't listen http: %v", err)
	}
	return nil
}
