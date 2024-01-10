package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/usrme/gobarchar"
)

var defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", timer(http.HandlerFunc(gobarchar.PresentBarChart)))

	log.Println("listening on:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func timer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)
		log.Println("completed in:", duration)
	})
}
