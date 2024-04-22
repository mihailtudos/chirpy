package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/mihailtudos/chirpy/internal/database"
	"github.com/mihailtudos/chirpy/pkg/utils"
)

func (a *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDSrt := r.PathValue("chirpID")
	id, err := strconv.Atoi(chirpIDSrt)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "incorrect id provided")
		return
	}

	chirp, err := a.db.GetChirp(id)
	if err != nil {
		if errors.Is(err, database.NotFound{}) {
			utils.RespondWithError(w, http.StatusNotFound, "couldn't get chirp")
			return
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, chirp)
}

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")

	dbChirps, err := a.db.GetChirps()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := dbChirps
	if authorId != "" {
		id, err := strconv.Atoi(authorId)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "author_id must be number")
		}

		chirps = filterChirpsByAuthor(dbChirps, id)
	}

	sortChirps(chirps, order)

	utils.RespondWithJSON(w, http.StatusOK, chirps)
}

func (a *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetTokenClaims(r.Header, a.jwtSecret)
	if err != nil {
		HandleTokenError(err, w)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	chirp, err := a.db.CreateChirp(cleaned, userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, database.Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorID: chirp.AuthorID,
	})
}

func (a *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDSrt := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDSrt)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "incorrect id provided")
		return
	}

	claims, err := a.GetTokenClaims(r.Header, a.jwtSecret)
	if err != nil {
		HandleTokenError(err, w)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = a.db.DeleteChirp(chirpID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotAuthorized) {
			utils.RespondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}

func filterChirpsByAuthor(chirps []database.Chirp, authorId int) []database.Chirp {
	filteredChirps := make([]database.Chirp, 0)
	for _, c := range chirps {
		if c.AuthorID == authorId {
			filteredChirps = append(filteredChirps, c)
		}
	}

	return filteredChirps
}

func sortChirps(chirps []database.Chirp, order string) {
	if order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
		return
	}
	
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
}