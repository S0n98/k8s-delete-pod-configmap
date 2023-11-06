package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)

func main() {
	var tlsCert string
	var tlsKey string
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Post("/validate", validate)

	flag.StringVar(&tlsCert, "tlsCertFile", "/etc/certs/cert.pem", "x509 Certificate for HTTPS.")
	flag.StringVar(&tlsKey, "tlsKeyFile", "/etc/certs/cert.pem", "x509 Key for HTTPS.")

	logger.Print("Starting server on port 3000...")
	http.ListenAndServeTLS(":3000", tlsCert, tlsKey, r)
}
