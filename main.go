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
var loggerErr = log.New(os.Stderr, "ERR: ", log.LstdFlags)

func main() {
	var tlsCert string
	var tlsKey string

	// init servemux
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Post("/validate", validate)

	// import cert and key from k8s secret
	flag.StringVar(&tlsCert, "tlsCertFile", "/etc/certs/tls.crt", "x509 Certificate for HTTPS.")
	flag.StringVar(&tlsKey, "tlsKeyFile", "/etc/certs/tls.key", "x509 Key for HTTPS.")

	logger.Print("Starting server on port 3443...")
	err := http.ListenAndServeTLS(":3443", tlsCert, tlsKey, r)
	if err != nil {
		loggerErr.Print(err)
	}
}
