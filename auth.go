package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihailtudos/chirpy/internal/database"
	"github.com/mihailtudos/chirpy/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const access_token = "chirpy-access"
const refresh_token = "chirpy-refresh"

var ErrAuthHeaderMissing error = errors.New("missing authrozation header")

type Claims struct {
	jwt.RegisteredClaims
}

func HandleTokenError(err error, w http.ResponseWriter) {
	fmt.Println(err)

	if errors.Is(err, jwt.ErrTokenInvalidClaims) {
		utils.RespondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	} else {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

// GetAuthTokenFromHeader returns the Bearer token from a given header
func (a *apiConfig) GetAuthTokenFromHeader(header http.Header) (string, error) {
	authorizationHeader := header.Get("Authorization")
	if authorizationHeader == "" {
		fmt.Println(ErrAuthHeaderMissing.Error())
		return "", ErrAuthHeaderMissing
	}

	token, _ := strings.CutPrefix(authorizationHeader, "Bearer ")

	return token, nil
}
func (a *apiConfig) GetTokenClaims(header http.Header, jwtSecret string) (Claims, error) {
	token, err := a.GetAuthTokenFromHeader(header)
	if err != nil {
		return Claims{}, err
	}

	claims := Claims{}
	_, err = jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
		jwt.WithLeeway(0*time.Second),
	)

	if err != nil {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}

	if claims.Issuer == refresh_token {
		ok, err := a.db.IsTokenRevoked(token)

		if err != nil {
			return Claims{}, err
		}

		if ok {
			return Claims{}, jwt.ErrTokenInvalidClaims
		}
	}

	return claims, nil
}

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	dj := json.NewDecoder(r.Body)
	var params parameters
	err := dj.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to parse the input")
		return
	}

	user, err := a.db.GetUserByEmail(params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "user not found")
		return
	}

	accessToken, err := a.issueJwtToken(access_token, user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	refreshToken, err := a.issueJwtToken(refresh_token, user.ID)
	if err != nil {
		fmt.Println("failed to sign\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (a *apiConfig) handlerRefreshTokens(w http.ResponseWriter, r *http.Request) {
	claims, err := a.GetTokenClaims(r.Header, a.jwtSecret)
	if err != nil {
		HandleTokenError(err, w)
		return
	}

	if claims.Issuer != refresh_token {
		utils.RespondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		fmt.Println("failed to cast string to int")
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	accessToken, err := a.issueJwtToken(access_token, userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, 200, struct {
		AccessToken string `json:"token"`
	}{
		AccessToken: accessToken,
	})
}

func (a *apiConfig) handlerRevokeRefreshTokens(w http.ResponseWriter, r *http.Request) {
	token, err := a.GetAuthTokenFromHeader(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = a.db.RevokeToken(database.Token(token))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}

func (a *apiConfig) issueJwtToken(issuer string, userId int) (string, error) {
	expiresIn := time.Hour * 24 * 60
	if issuer == access_token {
		expiresIn = time.Hour * 1
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   fmt.Sprint(userId),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(a.jwtSecret))
}
