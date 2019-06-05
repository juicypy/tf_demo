package models

import(
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type TFSession struct {
	Graph *tf.Graph
	Session *tf.Session
	Labels []string
}

type ClassifyResult struct {
	Labels   []LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

type ByProbability []LabelResult

func (a ByProbability) Len() int           { return len(a) }
func (a ByProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }
