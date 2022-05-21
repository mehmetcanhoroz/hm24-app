package handlers

import (
	"github.com/mehmetcanhoroz/hm24-app/services"
	"net/http"
)

type AnalyseHandler struct {
	AnalyserService services.IAnalyserService
}

func (h *AnalyseHandler) AnalysePage(w http.ResponseWriter, r *http.Request) {

}
