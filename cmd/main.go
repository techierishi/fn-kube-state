package main

import (
	"log"
	"net/http"
	"time"

	"fn-kube-state/pkg/handlers"
	"fn-kube-state/pkg/repository"
	"fn-kube-state/pkg/util"
)

func main() {

	dao := repository.NewDAO()
	srv := handlers.NewServer(dao)

	server := &http.Server{
		Addr:    ":8383",
		Handler: srv,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	util.Println("Starting server at port " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
