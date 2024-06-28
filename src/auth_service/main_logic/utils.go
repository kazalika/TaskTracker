package auth_service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"jwt_handlers"
	"mongo_handlers"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Generate JWT token for user
func GenerateJWTToken(username string) (string, error) {
	han := jwt_handlers.GetJWTHandlers()
	payload := jwt.MapClaims{
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)

	tokenString, err := token.SignedString(han.JwtPrivate)
	if err != nil {
		err = fmt.Errorf("error while signing token: %w", err)
		return "", err
	}
	return tokenString, nil
}

func CheckIfUserAuthenticated(r *http.Request, username *string) (code int, err error) {
	// Get token from Cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return http.StatusBadRequest, err
	}
	tokenString := cookie.Value

	han := jwt_handlers.GetJWTHandlers()
	payload := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &payload, func(token *jwt.Token) (interface{}, error) {
		return han.JwtPublic, nil
	})

	// Check if token is valid
	if err != nil || !token.Valid {
		return http.StatusBadRequest, errors.New("invalid jwt token")
	}

	// Check if token has neccessary information in payload
	if _, ok := payload["username"]; !ok {
		return http.StatusBadRequest, errors.New("invalid payload in jwt token")
	}
	*username = payload["username"].(string) // ?? is it ok?

	// Get user's relevant token
	relevantToken, err := mongo_handlers.GetUserToken(*username)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if token is relevant
	if relevantToken != tokenString {
		return http.StatusUnauthorized, errors.New("the token has expired")
	}

	return http.StatusOK, nil
}

// Hash given password by md5 + salt
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password + "SALT"))
	return fmt.Sprintf("%x", hash)
}

// Copy response `resp` to response writer `rw`
func CopyResponseToWriter(rw http.ResponseWriter, resp *http.Response) {
	rw.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	rw.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(rw, resp.Body)
	resp.Body.Close()
}
