package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("views"))
	http.Handle("/", fs)

	log.Println("Listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
