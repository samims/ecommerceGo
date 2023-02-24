package mocks

import (
	"mime/multipart"
)

type MockStorage struct {
	savedFile string
	err       error
}

func (m *MockStorage) Save(file multipart.File, header *multipart.FileHeader) (string, error) {
	return m.savedFile, m.err
}
