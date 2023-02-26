package utils

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface to string based JSON
func ToJSON(v interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

//// ToJSON serializes the given interface to a JSON string and returns it as a byte slice
//func ToJSON(v interface{}) ([]byte, error) {
//	buf := new(bytes.Buffer)
//	e := json.NewEncoder(buf)
//	err := e.Encode(v)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(v interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)

}
