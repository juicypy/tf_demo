package handlers

import (
	"bytes"
	"github.com/juicypy/tf_demo/models"
	"github.com/juicypy/tf_demo/usecases"
	"github.com/juicypy/tf_demo/utils"
	"io"
	"net/http"
	"strings"
)

type ImagesHandler interface {
	Recognize(w http.ResponseWriter, r *http.Request)
}

type imagesHandler struct {
	uc usecases.ImagesUsecase
}

func NewImagesHandler(u usecases.ImagesUsecase) ImagesHandler {
	return &imagesHandler{
		uc: u,
	}
}

func (s *imagesHandler) Recognize(w http.ResponseWriter, r *http.Request) {

	imageFile, header, err := r.FormFile("image")
	defer imageFile.Close()

	imageName := strings.Split(header.Filename, ".")
	if err != nil {
		utils.ResponseError(w, "Could not read image", http.StatusBadRequest)
		return
	}
	imageExt := imageName[len(imageName)-1]

	var imageBuffer bytes.Buffer
	_, err = io.Copy(&imageBuffer, imageFile)
	if err != nil {
		utils.ResponseError(w, "Could not read image", http.StatusBadRequest)
		return
	}

	res, err := s.uc.Recognize(&imageBuffer, imageExt)
	if err != nil {
		utils.ResponseError(w, "Could not run inference", http.StatusInternalServerError)
		return
	}

	utils.ResponseJSON(w, models.ClassifyResult{
		Labels: res,
	})
}
