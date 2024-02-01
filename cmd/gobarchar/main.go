package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/usrme/gobarchar"
)

var defaultPort = "8080"

type etagResponseWriter struct {
	http.ResponseWriter
	buf  bytes.Buffer
	hash hash.Hash
	w    io.Writer
}

func (e *etagResponseWriter) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", timer(etag(http.HandlerFunc(gobarchar.PresentBarChart))))

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

func etag(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &etagResponseWriter{
			ResponseWriter: w,
			buf:            bytes.Buffer{},
			hash:           sha1.New(),
		}
		ew.w = io.MultiWriter(&ew.buf, ew.hash)

		h.ServeHTTP(ew, r)

		etag := fmt.Sprintf("%x", ew.hash.Sum(nil))
		w.Header().Set("Etag", etag)

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
		} else {
			_, err := ew.buf.WriteTo(w)
			if err != nil {
				log.Println("unable to write HTTP response", err)
			}
		}
	})
}
