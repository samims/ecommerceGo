package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"product-images/configs"

	"github.com/sirupsen/logrus"
)

type Storage interface {
	Save(file multipart.File, header *multipart.FileHeader) (string, error)
}

type FileStorage struct {
	log      *logrus.Logger
	cfg      *configs.Config
	dir      string
	mediaURI string
}

func NewFileStorage(l *logrus.Logger, cfg *configs.Config) *FileStorage {
	return &FileStorage{
		log:      l,
		cfg:      cfg,
		dir:      cfg.ImageDIR,
		mediaURI: cfg.MediaURL,
	}
}

// Save saves a file to the file system using the provided multipart file and header.
// It generates a unique filename based on the current timestamp and file extension,
// creates a directory with the timestamp as its name, and saves the file to that directory.
// The function returns the complete filepath of the saved file, or an error if it fails to save the file.
func (s *FileStorage) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	// Generate a unique filename
	filename := header.Filename
	ext := filepath.Ext(filename)
	now := time.Now().Unix()
	filename = fmt.Sprintf("%d%s", now, ext)
	fileDir := filepath.Join(s.cfg.ImageDIR, fmt.Sprintf("%d", now))
	fileAbsPath := fmt.Sprintf("%s/%d/%s", s.mediaURI, now, filename)

	// Create the directory if it doesn't exist
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("unable to create directory: %s", err)
	}

	// Create the file
	filePathStr := filepath.Join(fileDir, filename)
	f, err := os.Create(filePathStr)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %s", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			s.log.Errorln("Log error ")
		}
	}(f)

	// Copy the contents of the request body to the file
	_, err = io.Copy(f, file)
	if err != nil {
		return "", fmt.Errorf("unable to save file: %s", err)
	}

	return fileAbsPath, nil
}
