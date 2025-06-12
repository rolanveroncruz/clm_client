package http01_challenge

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

const mapFilename = "./http01_challenge/http01.gob"

var challengeMap map[string]string

func init() {
	var loadErr error
	challengeMap, loadErr = LoadMap(mapFilename)
	if loadErr != nil {
		challengeMap = make(map[string]string)
	}

}

// PutChallenge accepts a token-authorization string pair to be saved for the challenge.
func PutChallenge(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Token      string `json:"token"`
		AuthString string `json:"AuthString"`
	}
	type Response struct {
		Status  int8   `json:"status"`
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json")
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error Decoding JSON", http.StatusBadRequest)
		return
	}
	var loadErr error
	challengeMap, loadErr = LoadMap(mapFilename)
	if loadErr != nil {
		http.Error(w, loadErr.Error(), http.StatusInternalServerError)
		return
	}
	challengeMap[requestData.Token] = requestData.AuthString

	saveErr := SaveMap(mapFilename, challengeMap)
	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusInternalServerError)
		return
	}
	response := Response{
		Status:  0,
		Message: "OK",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, writeErr := w.Write(jsonResponse)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

}

// GetChallenge returns the authorization string associated with the token
// in the path http://<your_domain>/.well-known/acme-challenge/<token>
func GetChallenge(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	token := params["token"]
	w.Header().Set("Content-Type", "text/plain")

	var loadErr error
	challengeMap, loadErr = LoadMap(mapFilename)
	if loadErr != nil {
		http.Error(w, loadErr.Error(), http.StatusInternalServerError)
		return
	}
	challenge := challengeMap[token]
	_, writeErr := w.Write([]byte(challenge))
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
	return

}

func SaveMap(filename string, data map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encodeErr := encoder.Encode(data)
	if encodeErr != nil {
		return encodeErr
	}
	return nil
}

func LoadMap(filename string) (map[string]string, error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return nil, openErr
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	decodeErr := decoder.Decode(&data)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return data, nil

}
