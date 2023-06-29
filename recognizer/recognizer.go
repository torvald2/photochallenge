package recognizer

import (
	"errors"

	fr "github.com/Kagami/go-face"
)

var NoFacesError = errors.New("No faces found")

type Tokenizer struct {
	rec       *fr.Recognizer
	tolerance float32
}

func (t Tokenizer) GetFaces(image []byte) (descriptor map[int][]float32, err error) {
	faces, err := t.rec.Recognize(image)
	if err != nil {
		return
	}
	if len(faces) == 0 {
		return descriptor, NoFacesError
	}
	for i, face := range faces {
		d := make([]float32, 128)
		desc := face.Descriptor
		for k, v := range desc {
			d[k] = v
		}
		descriptor[i] = d
	}
	return
}
func NewTokenizer(tolerance float32) (t Tokenizer, err error) {

	tmp_rec, err := fr.NewRecognizer("dnn_models")
	if err != nil {
		return
	}
	return Tokenizer{rec: tmp_rec, tolerance: tolerance}, nil
}
