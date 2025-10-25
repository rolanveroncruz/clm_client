package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"ph.certs.com/clm_client/auth"
	"ph.certs.com/clm_client/certs"
	http01 "ph.certs.com/clm_client/http01_challenge"
	"ph.certs.com/clm_client/middleware"
)

func main() {
	dotEnvErr := godotenv.Load()
	if dotEnvErr != nil {
		log.Fatal("Error loading .env file")
	}
	portStr := ":" + os.Getenv("CLM_CLIENT_PORT")
	muxRouter := mux.NewRouter()
	// POST /login accepts a username and password and returns a JWT token.
	// For now, we hardcode the username and password.
	muxRouter.Handle("/acme/login", middleware.LoggingMiddleware(middleware.CorsMiddleware(
		http.HandlerFunc(auth.Login)))).Methods("POST")

	// PUT /.well-known/acme-challenge/put-pair accepts a token-authorization string pair
	//to be saved for the HTTP-01 challenge.
	muxRouter.Handle("/acme/.well-known/acme-challenge/put-pair", middleware.LoggingMiddleware(
		middleware.JWTMiddleware(http.HandlerFunc(http01.PutChallenge)))).Methods("PUT")

	// GET acme-challenge responds to the HTTP-01 challenge.
	muxRouter.Handle("/acme/acme-challenge/{token}", middleware.LoggingMiddleware(
		http.HandlerFunc(http01.GetChallenge))).Methods("GET")

	muxRouter.Handle("/upload", middleware.LoggingMiddleware(
		middleware.JWTMiddleware(http.HandlerFunc(certs.UploadFileHandler)))).Methods("POST")

	muxRouter.Handle("/", middleware.LoggingMiddleware(http.FileServer(http.Dir("./static/")))).Methods("GET")

	muxRouter.NotFoundHandler = http.HandlerFunc(notFound)
	println("Listening on port " + portStr + "...")
	listenErr := http.ListenAndServe(portStr, muxRouter)
	if listenErr != nil {
		panic(listenErr)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("Not found %s %s\n", r.Method, r.URL.Path)
}
