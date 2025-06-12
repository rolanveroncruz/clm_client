package auth

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var SecretKey = []byte("mySecret-Key")

func Login(w http.ResponseWriter, r *http.Request) {
	catchAllErrorString := "Invalid Email or Password"
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type LoginResponse struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}

	var requestData LoginRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	email := requestData.Email
	password := requestData.Password
	if email == "admin@certs.com.ph" && password == "<PASSWORD>" {
		tokenString, err := createToken(email, "Administrator")
		response := LoginResponse{
			Name:  "Administrator",
			Email: email,
			Token: tokenString,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	} else {
		http.Error(w, catchAllErrorString, http.StatusBadRequest)
		return
	}
}
func createToken(userEmail string, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userEmail,
		"name":  userName,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
