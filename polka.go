package main

import (
	"encoding/json"
	"net/http"

	"github.com/mihailtudos/chirpy/pkg/utils"
)

const user_upgraded_event = "user.upgraded"

func (a *apiConfig) handlePolkaHook(w http.ResponseWriter, r *http.Request) {
	authKey, err := a.GetAuthTokenFromHeader(r.Header)
	if err != nil || authKey != a.polkaApiKey {
		utils.RespondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	jd := json.NewDecoder(r.Body)
	var param parameters
	err = jd.Decode(&param)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if param.Event != user_upgraded_event {
		utils.RespondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	if err := a.db.UpgradedUserToRedChirpy(param.Data.UserID); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
