package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
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

func (f *Files) saveFile(filenameWithExt, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", filenameWithExt, "path", path)
	// Create the directory if it doesn't exist
	dir := filepath.Join(f.cfg.ImageDIR)
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

//func (f *Files) UploadSingleFile(w http.ResponseWriter, r *http.Request) {
//	//file, header, err := ctx.Request.FormFile("image")
//	file, header, err := r.FormFile("image")
//	if err != nil {
//		http.Error(w, "Error fetching file from req", http.StatusInternalServerError)
//		return
//	}
//	unixStr := strconv.FormatInt(time.Now().Unix(), 10)
//	fileName := header.Filename
//	fileExt := filepath.Ext(fileName)
//	fileNameWithExt := fileName + unixStr + fileExt
//	f.saveFile(fileNameWithExt, f.cfg.ImageDIR, w, file)
//
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("Uploaded Successfully"))
//
//	response := struct {
//		Filepath string `json:"filepath"`
//	}{
//		Filepath: "filePathStr" + "x",
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(response)
//}

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
