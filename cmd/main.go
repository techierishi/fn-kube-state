package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"fn-kube-state/pkg/handlers"
	"fn-kube-state/pkg/repository"
)

func main() {

	ctx := context.Background()
	defer ctx.Done()

	dao := repository.NewDAO()
	srv := handlers.NewServer(dao, ctx)

	server := &http.Server{
		Addr:    ":8383",
		Handler: srv,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at port " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server started at port " + server.Addr)


}
