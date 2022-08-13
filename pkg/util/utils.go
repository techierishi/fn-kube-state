package util

import (
	"encoding/json"
	"log"
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
