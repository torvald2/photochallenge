package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type recognizer interface {
	GetFaces(image []byte) (descriptor map[int][]float32, err error)
	IsFacePresent(descriptor []float64, image []byte) (ok bool, err error)
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

type IsFacePresentRequest struct {
	Faces map[int][]float64 `json:"faces"`
	Image string            `json:"image"`
}

func (f FaceController) IsFacesPresent(c *gin.Context) {
	var payload IsFacePresentRequest

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err.Error()})
		return
	}
	img, err := base64.StdEncoding.DecodeString(payload.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err.Error()})
		return
	}
	for _, face := range payload.Faces {

		ok, err := f.recognizer.IsFacePresent(face, img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err.Error()})
			return
		}

		if !ok {
			c.JSON(http.StatusOK, map[string]interface{}{"status": "bad"})
			return
		}

	}
	c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})

}

func NewRouter(rg recognizer) *gin.Engine {
	router := gin.New()

	fc := FaceController{
		recognizer: rg,
	}

	router.POST("/api/v1/recognize", fc.GetFaceDescriptors)
	router.POST("/api/v1/compare", fc.IsFacesPresent)
	router.Use(static.Serve("/", static.LocalFile("./dist", false)))

	return router

}
