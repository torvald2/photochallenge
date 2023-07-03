package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type recognizer interface {
	GetFaces(image []byte) (descriptor map[int][]float32, err error)
}

type FaceController struct {
	recognizer recognizer
}
type RecognizeFacesRequest struct {
	Faces []string `json:"faces"`
}

type RecognizeFacesResponse struct {
	Descriptor map[int][]float32 `json:"descriptors"`
}

func (f FaceController) GetFaceDescriptors(c *gin.Context) {
	var payload RecognizeFacesRequest
	var resp RecognizeFacesResponse
	resp.Descriptor = make(map[int][]float32)

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err.Error()})
		return
	}
	for i, face := range payload.Faces {
		img, err := base64.StdEncoding.DecodeString(face)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err.Error()})
			return
		}
		faces, err := f.recognizer.GetFaces(img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
			return
		}
		for k, v := range faces {
			resp.Descriptor[k+i] = v
		}

	}
	c.JSON(http.StatusOK, resp)

}

func NewRouter(rg recognizer) *gin.Engine {
	router := gin.New()

	fc := FaceController{
		recognizer: rg,
	}

	router.POST("/api/v1/recognize", fc.GetFaceDescriptors)

	return router

}
