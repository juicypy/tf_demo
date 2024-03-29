package usecases

import (
	"bytes"
	"github.com/juicypy/tf_demo/models"
	"github.com/juicypy/tf_demo/repo"
	"sort"
)

type ImagesUsecase interface {
	Recognize(iBuffer *bytes.Buffer, iExt string) ([]models.LabelResult, error)
}

type imagesUsecase struct {
	repo   repo.ImagesRepo
	labels []string
}

func NewImagesUsecase(r repo.ImagesRepo, ls []string) ImagesUsecase {
	return &imagesUsecase{
		repo:   r,
		labels: ls,
	}
}

func (s *imagesUsecase) Recognize(iBuffer *bytes.Buffer, iExt string) ([]models.LabelResult, error) {
	output, err := s.repo.Recognize(iBuffer, iExt)
	if err != nil {
		return nil, err
	}
	res := s.findBestLabels(output[0].Value().([][]float32)[0])
	return res, nil
}

func (s *imagesUsecase) findBestLabels(probabilities []float32) []models.LabelResult {
	var resultLabels []models.LabelResult
	for i, p := range probabilities {
		if i >= len(s.labels) {
			break
		}
		resultLabels = append(resultLabels, models.LabelResult{Label: s.labels[i], Probability: p})
	}

	sort.Sort(models.ByProbability(resultLabels))
	return resultLabels[:5]
}
