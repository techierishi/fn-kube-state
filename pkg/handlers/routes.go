package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"fn-kube-state/pkg/repository"
)

type Server struct {
	router    *http.ServeMux
	db        repository.DAO
	muxRouter *mux.Router
}

func NewServer(doa repository.DAO) *Server {
	s := &Server{
		muxRouter: mux.NewRouter(),
		router:    http.NewServeMux(),
		db:        doa,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.muxRouter.ServeHTTP(w, r)
}

func (s *Server) routes() {

	s.muxRouter.HandleFunc("/pods", s.GetPods()).Methods("GET")
	s.muxRouter.HandleFunc("/services", s.GetDeployments()).Methods("GET")
	s.muxRouter.HandleFunc("/services/{appGroup}", s.GetDeploymentByGroup()).Methods("GET")
}

func (s *Server) ToJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	return e.Encode(data)
}
