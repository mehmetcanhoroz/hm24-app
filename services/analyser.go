package services

import (
	"github.com/mehmetcanhoroz/hm24-app/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
)

type IAnalyserService interface {
	GetHtmlContentOfURL(url string) string
	FindHtmlTitleOfURL(html string) string
	//DetermineHTMLVersion(html string) string
	//FindAllUrlsInPage(html string) []string
	//CountOfInternalUrls(urls []string) int
	//CountOfExternalUrls(urls []string) int
	//CountOfInaccessibleUrls(urls []string) int
	//FindHeadingsAndCounts(html string) map[string]int
	//DoesPageContainLoginForm(html string) bool
}

type AnalyseService struct {
}

func (s AnalyseService) FindHtmlTitleOfURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("While sending http request, an error occurred.", zap.Error(err))
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(resp.Body)

	htmlContent, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error("While parsing html content, an error occurred.", zap.Error(err))
		return ""
	}

	foundTitle, found := s.findDataOfHtmlElement(htmlContent, "title")

	if found {
		return foundTitle
	} else {
		return ""
	}
}

func (s AnalyseService) isThisElement(n *html.Node, elementType string) bool {
	return n.Type == html.ElementNode && n.Data == elementType
}

func (s AnalyseService) findDataOfHtmlElement(n *html.Node, elementType string) (string, bool) {
	if s.isThisElement(n, elementType) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := s.findDataOfHtmlElement(c, elementType)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func (s AnalyseService) GetHtmlContentOfURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("While sending http request, an error occurred.", zap.Error(err))
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(resp.Body)

	htmlOfURL, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("While reading html response, an error occurred.", zap.Error(err))
	}

	return string(htmlOfURL)
}
func NewAnalyseService() AnalyseService {
	return AnalyseService{}
}
