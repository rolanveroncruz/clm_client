package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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
	muxRouter.Handle("/login", middleware.LoggingMiddleware(middleware.CorsMiddleware(
		http.HandlerFunc(auth.Login)))).Methods("POST")

	muxRouter.Handle("/.well-known/acme-challenge/put-pair", middleware.LoggingMiddleware(
		middleware.JWTMiddleware(http.HandlerFunc(http01.PutChallenge)))).Methods("PUT")

	muxRouter.Handle("acme-challenge", middleware.LoggingMiddleware(
		http.HandlerFunc(http01.AcmeChallengeNoSlash))).Methods("GET")

	muxRouter.Handle("/acme-challenge", middleware.LoggingMiddleware(
		http.HandlerFunc(http01.AcmeChallengeSlashFrontOnly))).Methods("GET")

	muxRouter.Handle("/acme-challenge/", middleware.LoggingMiddleware(
		http.HandlerFunc(http01.AcmeChallengeSlashFrontBack))).Methods("GET")

	muxRouter.Handle("acme-challenge/", middleware.LoggingMiddleware(
		http.HandlerFunc(http01.AcmeChallengeSlashBackOnly))).Methods("GET")

	muxRouter.Handle("/acme-challenge/{token}", middleware.LoggingMiddleware(
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
