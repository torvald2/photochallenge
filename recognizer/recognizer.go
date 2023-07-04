package recognizer

import (
	"errors"
	"math"

	fr "github.com/Kagami/go-face"
	"gonum.org/v1/gonum/mat"
)

var NoFacesError = errors.New("No faces found")

type Tokenizer struct {
	rec       *fr.Recognizer
	tolerance float32
}

func (t Tokenizer) GetFaces(image []byte) (descriptor map[int][]float32, err error) {
	descriptor = make(map[int][]float32)
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

func (t Tokenizer) IsFacePresent(descriptor []float64, image []byte) (ok bool, err error) {
	userVector := mat.NewVecDense(128, descriptor)
	faces, err := t.rec.Recognize(image)
	if err != nil {
		return
	}
	for _, face := range faces {
		d := make([]float64, 128)
		desc := face.Descriptor
		for k, v := range desc {
			d[k] = float64(v)
		}
		vec := mat.NewVecDense(128, d)
		w := mat.NewVecDense(128, nil)
		w.SubVec(vec, userVector)
		dd := mat.Dot(w, w)
		dist := math.Sqrt(dd)
		if dist < 0.6 {
			return true, nil
		}

	}
	return false, nil

}
func NewTokenizer(tolerance float32) (t Tokenizer, err error) {

	tmp_rec, err := fr.NewRecognizer("dnn_models")
	if err != nil {
		return
	}
	return Tokenizer{rec: tmp_rec, tolerance: tolerance}, nil
}
