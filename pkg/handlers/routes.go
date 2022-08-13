package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	"fn-kube-state/pkg/repository"
)

type Server struct {
	router    *http.ServeMux
	db        repository.DAO
	muxRouter *mux.Router
	ctx       context.Context
}

func NewServer(doa repository.DAO, ctx context.Context) *Server {
	s := &Server{
		muxRouter: mux.NewRouter(),
		router:    http.NewServeMux(),
		db:        doa,
		ctx:       ctx,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.muxRouter.ServeHTTP(w, r)
}

func (s *Server) routes() {

	s.muxRouter.HandleFunc("/services", s.GetDeployments()).Methods("GET")
	s.muxRouter.HandleFunc("/services/{appGroup}", s.GetDeploymentByGroup()).Methods("GET")

	// Adding extra features out of interest
	s.muxRouter.HandleFunc("/stream", s.Stream()).Methods("GET")
	s.muxRouter.HandleFunc("/pods", s.GetPods()).Methods("GET")

	absPath, _ := filepath.Abs("./public")
	fmt.Println("absPath", absPath)
	s.muxRouter.PathPrefix("/").
		Handler(http.FileServer(http.Dir(absPath))).
		Methods("GET")
}

func (s *Server) ToJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/json")
	e := json.NewEncoder(w)
	return e.Encode(data)
}
