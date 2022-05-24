package services

import (
	"github.com/mehmetcanhoroz/hm24-app/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type IAnalyserService interface {
	GetHtmlContentOfURL(url string) string
	FindHtmlTitleOfURL(url string) string
	FindAllUrlsInPage(url string) []*html.Node
	FindAllUrlPathsInPage(elements []*html.Node) []string
	CountOfExternalUrlsInPage(elements []string) int
	DetermineHTMLVersion(htmlContent string) string
	//FindAllExternalPathsInPage(elements []*html.Node) []string
	//FindAllInternalPathsInPage(elements []*html.Node) []string
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

func (s AnalyseService) CountOfExternalUrlsInPage(urls []string) int {
	var count int
	externalIdentifiers := []string{
		"http://",
		"https://",
	}
	for _, url := range urls {
		for _, identifier := range externalIdentifiers {
			if strings.HasPrefix(strings.TrimSpace(url), identifier) {
				count++
				break
			}
		}
	}
	return count
}
func (s AnalyseService) DetermineHTMLVersion(htmlContent string) string {
	var htmlVersions = make(map[string]string)
	var version = "UNKNOWN"

	htmlVersions["HTML 4.01 Strict"] = `"-//W3C//DTD HTML 4.01//EN"`
	htmlVersions["HTML 4.01 Transitional"] = `"-//W3C//DTD HTML 4.01 Transitional//EN"`
	htmlVersions["HTML 4.01 Frameset"] = `"-//W3C//DTD HTML 4.01 Frameset//EN"`
	htmlVersions["HTML 4.01"] = `DTD HTML 4.01`
	htmlVersions["XHTML 1.0 Strict"] = `"-//W3C//DTD XHTML 1.0 Strict//EN"`
	htmlVersions["XHTML 1.0 Transitional"] = `"-//W3C//DTD XHTML 1.0 Transitional//EN"`
	htmlVersions["XHTML 1.0 Frameset"] = `"-//W3C//DTD XHTML 1.0 Frameset//EN"`
	htmlVersions["XHTML 1.0 Frameset"] = `"DTD XHTML 1.0 Frameset`
	htmlVersions["XHTML 1.1"] = `"-//W3C//DTD XHTML 1.1//EN"`
	htmlVersions["HTML 5"] = `<!DOCTYPE html>`

	for doctype, matcher := range htmlVersions {
		match := strings.Contains(htmlContent, matcher)

		if match == true {
			version = doctype
		}
	}

	return version
}

func NewAnalyseService() AnalyseService {
	return AnalyseService{}
}
