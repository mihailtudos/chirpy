package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func ReplaceProfaneWords(text string) string {
	words := strings.Fields(text)
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		if InSlice(word, profaneWords) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func InSlice(word string, words []string) bool {
	for _, w := range words {
		if strings.EqualFold(w, word) {
			return true
		}
	}

	return false
}
