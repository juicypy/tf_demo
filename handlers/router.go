package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/juicypy/tf_demo/models"
	"github.com/juicypy/tf_demo/repo"
	"github.com/juicypy/tf_demo/usecases"
	"github.com/juicypy/tf_demo/utils"
	"net/http"
)

func NewRouter(ctx context.Context)*mux.Router{
	router := mux.NewRouter()
	addRoutes(router, ctx)
	return router
}

func addRoutes(router *mux.Router, ctx context.Context){
	tfSess, ok := ctx.Value(utils.TFSessionCtx).(*models.TFSession)
	if !ok{
		panic("asserting session error")
	}

	ir := repo.NewImagesRepo(tfSess)
	iuc := usecases.NewImagesUsecase(ir, tfSess.Labels)
	ih := NewImagesHandler(iuc)

	router.HandleFunc("/recognize/image", ih.Recognize).Methods(http.MethodPost)
}
