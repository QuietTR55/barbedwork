package utilities

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadDir = "./uploads"

func SaveImage(w http.ResponseWriter, r *http.Request) (string, error) {

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return "", err
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("Error getting file from form:", err)
		http.Error(w, "Unable to get the image", http.StatusBadRequest)
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	//Decode the image
	file.Seek(0, 0)
	_, _, err = image.Decode(file)
	if err != nil {
		log.Println("Error decoding image:", err)
		http.Error(w, "Invalid image format", http.StatusBadRequest)
		return "", err
	}

	file.Seek(0, 0)

	if os.MkdirAll(uploadDir, os.ModePerm) != nil {
		log.Println("Error creating upload directory:", err)
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return "", err
	}

	out, err := os.Create(uploadDir + "/" + filename)
	if err != nil {
		log.Println("Error creating file:", err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return "", err
	}
	defer out.Close()

	file.Seek(0, 0)
	if _, err := out.ReadFrom(file); err != nil {
		log.Println("Error saving file:", err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return "", err
	}

	filePath := "/uploads/" + filename // Changed to a URL path
	return filePath, nil
}
