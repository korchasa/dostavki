package web

import (
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) handleStatic() http.HandlerFunc {
	var (
		staticPath = "./web/static"
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var path string
		if "/" == r.URL.Path {
			path = "/index.html"
		} else {
			path = r.URL.Path
		}

		path, err := filepath.Abs(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		path = filepath.Join(staticPath, path)

		_, err = os.Stat(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
	}
}
