package certs

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadFileHandler handles file uploads in the route "/upload". It accepts the file
// and writes it to file path described by the environment variable "CERT_FILE_PATH".
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the multipart form
	// Max 10 MB file size.
	// A certificate file is about 4k in size; private-key is 2k; ca cert is 2k. So 10MB is more than enough.
	uploadDir := os.Getenv("CERT_FILE_PATH")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		fmt.Printf("%s does not exist.", uploadDir)
		createErr := os.Mkdir(uploadDir, 0755)
		if createErr != nil {
			fmt.Printf("Error creating upload directory: %s\n", createErr)
		}
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB max..
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing multipart form: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Retrieve the file from the form data
	file, handler, err := r.FormFile("myFile") // "myFile" is the key for the file in the form
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file from form: %v", err), http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(file)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. Create a new file on the server to save the uploaded data
	dstPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating file on server: %v", err), http.StatusInternalServerError)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(dst)

	// 4. Copy the uploaded file's content to the new file on the server
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error copying file content: %v", err), http.StatusInternalServerError)
		return
	}

	_, printErr := fmt.Fprintf(w, "Successfully Uploaded File: %s\n", handler.Filename)
	if printErr != nil {
		return
	}
	// TODO: Restart NGINX: sudo systemctl restart nginx
}
