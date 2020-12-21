package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// MetaData is meta data struct
type MetaData = toml.MetaData

// Encoder is encoder struct
type Encoder struct {
	encoder *toml.Encoder
}

// NewEncoder returns a new encoder instance
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		encoder: toml.NewEncoder(w),
	}
}

// Encode encodes interface to toml
func (enc *Encoder) Encode(v interface{}) error {
	return enc.encoder.Encode(v)
}

// Decode decodes toml data to interface
func Decode(data string, v interface{}) (MetaData, error) {
	return toml.Decode(data, v)
}

// DecodeFile decodes toml file
func DecodeFile(fpath string, v interface{}) (MetaData, error) {
	return toml.DecodeFile(fpath, v)
}
