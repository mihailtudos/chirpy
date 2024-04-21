package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mihailtudos/chirpy/internal/database"
	"github.com/mihailtudos/chirpy/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func (a *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	j := json.NewDecoder(r.Body)
	err := j.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if !utils.IsPasswordValid(params.Password) {
		utils.RespondWithError(w, http.StatusBadRequest, "password must be at least 8 characters long")
		return
	}

	hashedBytePass, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	user, err := a.db.CreateUser(database.User{Email: params.Email, Password: string(hashedBytePass)})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}{
		Email: user.Email,
		ID:    user.ID,
	})
}

func (a *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetTokenClaims(r.Header, a.jwtSecret)
	if err != nil {
		HandleTokenError(err, w)
		return
	}

	if claims.Issuer == refresh_token {
		utils.RespondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	jd := json.NewDecoder(r.Body)
	params := parameters{}
	if err := jd.Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, err := a.db.UpdateUser(userId, params.Email, params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}{
		Email: user.Email,
		ID:    user.ID,
	})
}
