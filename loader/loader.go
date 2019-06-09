package loader

import (
	"bufio"
	"github.com/juicypy/tf_demo/models"
	"io/ioutil"
	"log"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Loader struct {
	ModelPath  string
	LabelsPath string
}

func (s *Loader) NewSession() (*models.TFSession, error) {

	model, err := ioutil.ReadFile(s.ModelPath)
	if err != nil {
		return nil, err
	}
	graphModel := tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return nil, err
	}

	labelsFile, err := os.Open(s.LabelsPath)
	if err != nil {
		return nil, err
	}
	defer labelsFile.Close()

	scanner := bufio.NewScanner(labelsFile)

	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sessionModel, err := tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	sess := &models.TFSession{
		Graph:   graphModel,
		Session: sessionModel,
		Labels:  labels,
	}

	return sess, nil
}
