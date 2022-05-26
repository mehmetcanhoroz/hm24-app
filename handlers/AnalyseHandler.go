package handlers

import (
	"encoding/json"
	"github.com/mehmetcanhoroz/hm24-app/logger"
	"github.com/mehmetcanhoroz/hm24-app/services"
	"github.com/mehmetcanhoroz/hm24-app/utils/rest_utils"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

type AnalyseHandler struct {
	AnalyserService services.IAnalyserService
}

func (h *AnalyseHandler) GetHtmlContentOfURL(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent := h.AnalyserService.SendHttpRequest(url)
	htmlContent := h.AnalyserService.GetHtmlContentOfURL(requestContent)

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

func (h *AnalyseHandler) DetermineHTMLVersion(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent := h.AnalyserService.SendHttpRequest(url)

	htmlContent := h.AnalyserService.GetHtmlContentOfURL(requestContent)

	err := rest_utils.PrepareApiResponseAsJson(w)
	if err != nil {
		return
	}

	version := h.AnalyserService.DetermineHTMLVersion(htmlContent)

	response := rest_utils.NewApiResponse(200, version, "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *AnalyseHandler) FindHtmlTitleOfURL(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent := h.AnalyserService.SendHttpRequest(url)
	htmlContent := h.AnalyserService.FindHtmlTitleOfURL(requestContent)

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

func (h *AnalyseHandler) GetListOfLinkElements(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent := h.AnalyserService.SendHttpRequest(url)

	requestContentX, _ := html.Parse(requestContent)

	htmlContent := h.AnalyserService.FindAllXElementInPage(requestContentX, "a")

	err := rest_utils.PrepareApiResponseAsJson(w)
	if err != nil {
		return
	}

	links := h.AnalyserService.FindAllUrlPathsInPage(htmlContent)
	externalUrlCount := h.AnalyserService.CountOfExternalUrlsInPage(links)

	countOfAccessibleLinks := h.AnalyserService.CountOfAccessibleUrls(links, url)

	responseWithCount := map[string]interface{}{
		"total_count":             len(links),
		"external_count":          externalUrlCount,
		"internal_count":          len(links) - externalUrlCount,
		"links":                   links,
		"accessible_link_count":   countOfAccessibleLinks,
		"inaccessible_link_count": len(links) - countOfAccessibleLinks,
	}

	response := rest_utils.NewApiResponse(200, responseWithCount, "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *AnalyseHandler) GetCountOfHXElements(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent0 := h.AnalyserService.SendHttpRequest(url)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(requestContent0)
	//defer func(Body io.ReadCloser) {
	//	if closeRequestNow {
	//		err := Body.Close()
	//		if err != nil {
	//			logger.Error("While closing http request, an error occurred.", zap.Error(err))
	//		}
	//	}
	//}(htmlContent)

	requestContentX, err := html.Parse(requestContent0)
	if err != nil {
		logger.Error("While parsing html content, an error occurred.", zap.Error(err))
		return
	}

	requestContent1 := requestContentX
	requestContent2 := requestContentX
	requestContent3 := requestContentX
	requestContent4 := requestContentX
	requestContent5 := requestContentX
	requestContent6 := requestContentX
	h1HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent1, "h1")
	h2HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent2, "h2")
	h3HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent3, "h3")
	h4HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent4, "h4")
	h5HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent5, "h5")
	h6HtmlContent := h.AnalyserService.FindAllXElementInPage(requestContent6, "h6")

	err = rest_utils.PrepareApiResponseAsJson(w)
	if err != nil {
		return
	}

	responseWithCount := map[string]interface{}{
		"h1_count": len(h1HtmlContent),
		"h2_count": len(h2HtmlContent),
		"h3_count": len(h3HtmlContent),
		"h4_count": len(h4HtmlContent),
		"h5_count": len(h5HtmlContent),
		"h6_count": len(h6HtmlContent),
	}

	response := rest_utils.NewApiResponse(200, responseWithCount, "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (h *AnalyseHandler) IsThereLoginForm(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams["url"][0]

	requestContent := h.AnalyserService.SendHttpRequest(url)

	requestContentX, _ := html.Parse(requestContent)
	requestContent1 := requestContentX
	requestContent2 := requestContentX

	err := rest_utils.PrepareApiResponseAsJson(w)

	htmlInputContent := h.AnalyserService.FindAllXElementInPage(requestContent1, "input")
	htmlFormContent := h.AnalyserService.FindAllXElementInPage(requestContent2, "form")

	println(htmlFormContent)
	println(htmlInputContent)

	if err != nil {
		return
	}

	response := rest_utils.NewApiResponse(200, h.AnalyserService.IsThereLoginForm(htmlInputContent, htmlFormContent), "")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func NewAnalyseHandler(analyseService services.IAnalyserService) AnalyseHandler {
	return AnalyseHandler{AnalyserService: analyseService}
}
