package repo

import (
	"bytes"
	"github.com/juicypy/tf_demo/models"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"

	"log"
)

const (
	pngExt           = "png"
	rgbChannelsCount = 3
	inputOperation   = "input"
	outputOperation  = "output"
)

type ImagesRepo interface {
	Recognize(iBuffer *bytes.Buffer, iExt string) ([]*tf.Tensor, error)
}

type imagesRepo struct {
	sess *models.TFSession
}

func NewImagesRepo(s *models.TFSession) ImagesRepo {
	return &imagesRepo{
		sess: s,
	}
}

type transImage struct {
	Graph         *tf.Graph
	Input, Output tf.Output
}

func (s *imagesRepo) Recognize(iBuffer *bytes.Buffer, iExt string) ([]*tf.Tensor, error) {
	tensor, err := s.tensorFromImage(iBuffer, iExt)
	if err != nil {
		return nil, err
	}

	// Run inference
	output, err := s.sess.Session.Run(
		map[tf.Output]*tf.Tensor{
			s.sess.Graph.Operation(inputOperation).Output(0): tensor,
		},
		[]tf.Output{
			s.sess.Graph.Operation(outputOperation).Output(0),
		},
		nil,
	)

	return output, nil
}

func (s *imagesRepo) tensorFromImage(imageBuffer *bytes.Buffer, imageFormat string) (*tf.Tensor, error) {
	tensor, err := tf.NewTensor(imageBuffer.String())
	if err != nil {
		return nil, err
	}
	tsImage, err := s.transformImageGraph(imageFormat)
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(tsImage.Graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{tsImage.Input: tensor},
		[]tf.Output{tsImage.Output},
		nil,
	)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func (s *imagesRepo) transformImageGraph(imageFormat string) (transImage, error) {
	const (
		H, W  = 224, 224
		Mean  = float32(117)
		Scale = float32(1)
	)

	scope := op.NewScope()
	input := op.Placeholder(scope, tf.String)
	// Decode PNG or JPEG
	var decode tf.Output
	if imageFormat == pngExt {
		decode = op.DecodePng(scope, input, op.DecodePngChannels(rgbChannelsCount))
	} else {
		decode = op.DecodeJpeg(scope, input, op.DecodeJpegChannels(rgbChannelsCount))
	}

	// Resize to 224x224 with bilinear interpolation
	resized := op.ResizeBilinear(scope,
		// Create a batch containing a single image
		op.ExpandDims(scope,
			// Use decoded pixel values
			op.Cast(scope, decode, tf.Float),
			op.Const(scope.SubScope("make_batch"), int32(0))),
		op.Const(scope.SubScope("size"), []int32{H, W}),
	)

	// Div and Sub perform (value-Mean)/Scale for each pixel
	output := op.Div(scope,
		op.Sub(scope, resized, op.Const(scope.SubScope("mean"), Mean)),
		op.Const(scope.SubScope("scale"), Scale),
	)
	graph, err := scope.Finalize()

	res := transImage{
		graph,
		input,
		output,
	}
	return res, err
}
