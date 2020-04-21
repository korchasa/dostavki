package web

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleStatic())
	s.router.HandleFunc("/stores/{lat}/{long}", s.handleStores())
}