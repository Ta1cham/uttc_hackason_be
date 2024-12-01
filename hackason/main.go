package main

import (
	"encoding/json"
	"net/http"
)

type responseMessage struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(responseMessage{
		Message: "Hello " + name + "-san.",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func main() {
	http.HandleFunc("/hello", handler)
	http.ListenAndServe(":8000", nil)
}
