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

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(v interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)

}
