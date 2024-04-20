package main

import (
	"encoding/json"
	"net/http"

	"github.com/mihailtudos/chirpy/internal/database"
	"github.com/mihailtudos/chirpy/pkg/utils"
)

func (a *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	params := parameters{}
	j := json.NewDecoder(r.Body)
	err := j.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := a.db.CreateUser(database.User{Email: params.Email})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, user)
}
