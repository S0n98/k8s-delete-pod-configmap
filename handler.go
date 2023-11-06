package main

import (
	"io"
	"net/http"
)

func validate(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error

	bodyBytes, err = io.ReadAll(r.Body)

	if err != nil {
		logger.Print("Reading body failed")
	}

	w.Write(bodyBytes)
}
