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
		loggerErr.Print("Reading body failed")
	}

	logger.Print(string(bodyBytes))
	w.Write(bodyBytes)
}
