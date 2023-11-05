package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct {
		Msg string `json:"msg"`
	}{
		Msg: "everything is working like a charm",
	})
}
