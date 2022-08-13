package handlers

import (
	"encoding/json"
	"fmt"
	"fn-kube-state/pkg/models"
	"fn-kube-state/pkg/util"
	"log"
	"net/http"
	"time"

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

func (s *Server) Stream() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		client := &models.Client{Name: r.RemoteAddr, Events: make(chan *models.SseMessage, 10)}
		go func() {
			s.db.NewKubeQuery().Watch(s.ctx, client)
		}()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(10 * time.Second)
		select {
		case ev := <-client.Events:
			a, _ := json.Marshal(ev)

			fmt.Fprintf(w, "data: %v\n\n", string(a))
			fmt.Printf("data: %v\n", string(a))
		case <-timeout:
			fmt.Fprintf(w, ": Nothing to sent\n\n")
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
