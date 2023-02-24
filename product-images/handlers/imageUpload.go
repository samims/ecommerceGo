package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"product-images/configs"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	Save(file multipart.File, header *multipart.FileHeader) (string, error)
}

type Files struct {
	log   *logrus.Logger
	store Storage
	cfg   *configs.Config
}

func NewFiles(s Storage, l *logrus.Logger, cfg *configs.Config) *Files {
	return &Files{
		log:   l,
		store: s,
		cfg:   cfg,
	}
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(w, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) UploadMultipart(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error fetching file from req", http.StatusInternalServerError)
		return
	}

	// Save the file to storage
	filePath, err := f.store.Save(file, header)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Return the file path to the client
	response := struct {
		Filepath string `json:"filepath"`
	}{
		Filepath: filePath,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	buf, err := os.ReadFile("sid.png")

	if err != nil {

		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", `attachment;filename="sid.png"`)

	w.Write(buf)
}

func (f *Files) ServeImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unixTime := vars["unixTime"]
	fileName := vars["fileName"]
	imagePath := fmt.Sprintf("%s/%s/%s", f.cfg.ImageDIR, unixTime, fileName)

	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			f.log.Errorf("Closing file failed after reading from serve image for path %s", imagePath)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}
