package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func PrintJSON(key string, data interface{}) {
	if os.Getenv("DEBUG") == "true" {
		a, _ := json.Marshal(data)
		log.Println(key, string(a))
	}
}

func Println(args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		log.Println(args)
	}
}

func PrepareHeaderForSSE(w http.ResponseWriter) {
	// prepare the header
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func WriteData(w http.ResponseWriter, messageChan chan string) (int, error) {
	// set data into response writer
	return fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
}
