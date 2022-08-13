package handlers

import (
	"fn-kube-state/pkg/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) GetDeployments() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		deps, err := s.db.NewKubeQuery().GetDeploymentByGroup(r.Context(), "default", "")
		if err != nil {
			log.Fatal(err)
		}
		s.ToJSON(w, deps)

	}
}

func (s *Server) GetDeploymentByGroup() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		appGroup := vars["appGroup"]
		util.Println("appGroup", appGroup)
		deps, err := s.db.NewKubeQuery().GetDeploymentByGroup(r.Context(), "default", appGroup)
		if err != nil {
			log.Fatal(err)
		}
		s.ToJSON(w, deps)

	}
}

func (s *Server) GetPods() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		pods, err := s.db.NewKubeQuery().GetPods(r.Context())
		if err != nil {
			log.Fatal(err)
		}
		s.ToJSON(w, pods)

	}
}
