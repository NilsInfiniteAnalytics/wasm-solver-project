package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Setting up server...")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Println("Listening on http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
