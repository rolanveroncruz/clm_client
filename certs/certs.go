package certs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads" // Directory to save uploaded uploads

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the multipart form
	// Max 10 MB file size.
	// A certificate file is about 4k in size; private-key is 2k; ca cert is 2k. So 10MB is more than enough.
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
	defer file.Close()

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
	defer dst.Close()

	// 4. Copy the uploaded file's content to the new file on the server
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error copying file content: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File: %s\n", handler.Filename)
}
