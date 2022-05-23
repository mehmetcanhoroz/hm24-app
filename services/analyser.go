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
	FindHtmlTitleOfURL(url string) string
	FindAllUrlsInPage(url string) []*html.Node
	FindAllUrlPathsInPage(elements []*html.Node) []string
	//FindAllExternalPathsInPage(elements []*html.Node) []string
	//FindAllInternalPathsInPage(elements []*html.Node) []string
	//DetermineHTMLVersion(html string) string
	//
	//CountOfInternalUrls(urls []string) int
	//CountOfExternalUrls(urls []string) int
	//CountOfInaccessibleUrls(urls []string) int
	//FindHeadingsAndCounts(html string) map[string]int
	//DoesPageContainLoginForm(html string) bool
}

type AnalyseService struct {
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

func (s AnalyseService) getListOfTypeHtmlElements(n *html.Node, elementType string, listOfElements []*html.Node) []*html.Node {

	if s.isThisElement(n, elementType) {
		listOfElements = append(listOfElements, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		listOfElements = s.getListOfTypeHtmlElements(c, elementType, listOfElements)
	}

	return listOfElements
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

func (s AnalyseService) FindAllUrlPathsInPage(elements []*html.Node) []string {
	var links []string
	for _, node := range elements {
		links = append(links, node.Attr[0].Val)
	}
	return links
}

func (s AnalyseService) FindAllUrlsInPage(url string) []*html.Node {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("While sending http request, an error occurred.", zap.Error(err))
		return nil
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
		return nil
	}

	var tempList []*html.Node

	listLinks := s.getListOfTypeHtmlElements(htmlContent, "a", tempList)

	return listLinks
}

func NewAnalyseService() AnalyseService {
	return AnalyseService{}
}
