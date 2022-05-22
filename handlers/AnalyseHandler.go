package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mehmetcanhoroz/hm24-app/services"
	"github.com/mehmetcanhoroz/hm24-app/utils/rest_utils"
	"net/http"
)

type AnalyseHandler struct {
	AnalyserService services.IAnalyserService
}

func (h *AnalyseHandler) GetHtmlContentOfURL(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	htmlContent := h.AnalyserService.GetHtmlContentOfURL(fmt.Sprintf("%s", url))

	err := rest_utils.PrepareApiResponseAsJson(w)
	if err != nil {
		return
	}

	response := rest_utils.NewApiResponse(200, htmlContent, "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *AnalyseHandler) FindHtmlTitleOfURL(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	htmlContent := h.AnalyserService.FindHtmlTitleOfURL(fmt.Sprintf("%s", url))

	err := rest_utils.PrepareApiResponseAsJson(w)
	if err != nil {
		return
	}

	response := rest_utils.NewApiResponse(200, htmlContent, "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func NewAnalyseHandler(analyseService services.IAnalyserService) AnalyseHandler {
	return AnalyseHandler{AnalyserService: analyseService}
}
