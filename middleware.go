package main

import (
	"encoding/json"
	"net/http"

	"github.com/VividCortex/siesta"
)

type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

func jsonMarshaler(c siesta.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	errorString := ""
	contextErr := c.Get(errorKey)
	if contextErr != nil {
		errorString = contextErr.(string)
	}

	response := apiResponse{
		Data:  c.Get(responseKey),
		Error: errorString,
	}

	json.NewEncoder(w).Encode(response)
}
