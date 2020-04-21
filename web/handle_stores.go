package web

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
)

func (s *Server) handleStores() http.HandlerFunc {
	tpl, err := template.ParseFiles("./web/templates/stores.gohtml")
	if err != nil {
		log.Panicf("can't parse index template: %v", err)
	}

	onError := func(w http.ResponseWriter, err error) {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		lat, err := strconv.ParseFloat(params["lat"], 32)
		if err != nil {
			onError(w, err)
			return
		}
		long, err := strconv.ParseFloat(params["long"], 32)
		if err != nil {
			onError(w, err)
			return
		}
		scheds, err := s.zakazClient.Schedule(lat, long)
		if err != nil {
			onError(w, err)
			return
		}

		if err := tpl.Execute(w, scheds); err != nil {
			onError(w, err)
		}
	}
}
