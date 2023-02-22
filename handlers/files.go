package handlers

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/samims/ecommerceGo/configs"
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

func NewFiles(l *logrus.Logger, cfg *configs.Config) *Files {
	return &Files{
		log: l,
		//store: s,
		cfg: cfg,
	}
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(w, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) saveFile(filenameWithExt, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", filenameWithExt, "path", path)
	// Create the directory if it doesn't exist
	dir := filepath.Join(f.cfg.FileDir)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		f.log.Error("Unable to create directory", "dir", dir, "error", err)
		http.Error(w, "Unable to create directory", http.StatusInternalServerError)
		return
	}

	// Create the file
	file, err := os.Create(filepath.Join(dir, filenameWithExt))
	if err != nil {
		f.log.Error("Unable to create file", "path", path, "error", err)
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			f.log.Error("Unable to close file", path)
		}
	}(file)

	// Copy the contents of the request body to the file
	_, err = io.Copy(file, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	f.log.Info("File saved successfully")
}

func (f *Files) UploadSingleFile(w http.ResponseWriter, r *http.Request) {
	//file, header, err := ctx.Request.FormFile("image")
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error fetching file from req", http.StatusInternalServerError)
		return
	}
	unixStr := strconv.FormatInt(time.Now().Unix(), 10)
	fileName := header.Filename
	fileExt := filepath.Ext(fileName)
	fileNameWithExt := fileName + unixStr + fileExt
	f.saveFile(fileNameWithExt, f.cfg.FileDir, w, file)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Uploaded Successfully"))

	response := struct {
		Filepath string `json:"filepath"`
	}{
		Filepath: "filePathStr" + "x",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
