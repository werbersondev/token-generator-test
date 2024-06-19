package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
	LivenessHandler               http.HandlerFunc
	RequestTokenGenerationHandler http.HandlerFunc
}

func New(service RequestTokenGenerationUseCase) *API {
	api := API{
		LivenessHandler:               LivenessHandler(),
		RequestTokenGenerationHandler: RequestTokenGenerationHandler(service),
	}

	return &api
}

func (a *API) Routes(router *chi.Mux) {
	router.HandleFunc("/liveness", a.LivenessHandler)
	router.HandleFunc("/generate_token", a.RequestTokenGenerationHandler)
}

func LivenessHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
