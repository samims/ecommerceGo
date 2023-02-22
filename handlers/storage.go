package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalDiskStorage implements the Storage interface for local disk storage.
type LocalDiskStorage struct {
	path   string
	prefix string
}

// NewLocalDiskStorage returns a new instance of LocalDiskStorage.
func NewLocalDiskStorage(path, prefix string) *LocalDiskStorage {
	return &LocalDiskStorage{
		path:   path,
		prefix: prefix,
	}
}

// Save saves a file to the local disk storage.
func (s *LocalDiskStorage) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(s.path, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %s", err)
	}

	// Generate a unique file name for the file.
	fileName := GenerateFileName(s.prefix, header.Filename)

	// Create the destination file.
	dst, err := os.Create(filepath.Join(s.path, fileName))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		// TODO: add log by dependency file injection
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(dst)

	// Copy the contents of the uploaded file to the destination file.
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %s", err)
	}

	return fileName, nil
}

// GenerateFileName generates a unique file name with the specified prefix and original file name.
func GenerateFileName(prefix, filename string) string {
	ext := filepath.Ext(filename)
	return fmt.Sprintf("%s%s%s", prefix, RandomString(10), ext)
}

// RandomString generates a random string with the specified length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
