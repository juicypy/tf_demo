package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/juicypy/tf_demo/handlers"
	"github.com/juicypy/tf_demo/loader"
	"github.com/juicypy/tf_demo/utils"
	"log"
	"net/http"
)

func main() {

	modelPath := flag.String("mp", "tf_model.pb", "model_path")
	labelsPath := flag.String("lp", "label_strings.txt", "labels_path")
	flag.Parse()

	l := loader.Loader{
		ModelPath: *modelPath,
		LabelsPath: *labelsPath,
	}

	sess, err := l.NewSession()
	if err != nil{
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.TFSessionCtx, sess)

	router := handlers.NewRouter(ctx)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

