package services

import (
	"github.com/mehmetcanhoroz/hm24-app/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type IAnalyserService interface {
	SendHttpRequest(url string) io.ReadCloser
	GetHtmlContentOfURL(htmlContent io.ReadCloser) string
	FindHtmlTitleOfURL(htmlContent io.ReadCloser) string
	FindAllXElementInPage(gotHtmlContent *html.Node, elementType string) []*html.Node
	FindAllUrlPathsInPage(elements []*html.Node) []string
	CountOfExternalUrlsInPage(elements []string) int
	DetermineHTMLVersion(htmlContent string) string
	IsThereLoginForm(inputs []*html.Node, forms []*html.Node) bool
	CountOfAccessibleUrls(urls []string, baseUrl string) int
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

func (s AnalyseService) SendHttpRequest(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("While sending http request, an error occurred.", zap.Error(err))
		return nil
	}

	/*defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(resp.Body)*/

	return resp.Body
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

func (s AnalyseService) GetHtmlContentOfURL(htmlContent io.ReadCloser) string {

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(htmlContent)

	htmlOfURL, err := ioutil.ReadAll(htmlContent)
	if err != nil {
		logger.Error("While reading html response, an error occurred.", zap.Error(err))
	}

	return string(htmlOfURL)
}

func (s AnalyseService) FindHtmlTitleOfURL(htmlContent io.ReadCloser) string {

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("While closing http request, an error occurred.", zap.Error(err))
		}
	}(htmlContent)

	gotHtmlContent, err := html.Parse(htmlContent)
	if err != nil {
		logger.Error("While parsing html content, an error occurred.", zap.Error(err))
		return ""
	}

	foundTitle, found := s.findDataOfHtmlElement(gotHtmlContent, "title")

	if found {
		return foundTitle
	} else {
		return ""
	}
}

func (s AnalyseService) FindAllUrlPathsInPage(elements []*html.Node) []string {
	var links []string
	for _, node := range elements {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" {
				links = append(links, attribute.Val)
			}
		}
	}
	return links
}

func (s AnalyseService) FindAllXElementInPage(gotHtmlContent *html.Node, elementType string) []*html.Node {

	var tempList []*html.Node

	listLinks := s.getListOfTypeHtmlElements(gotHtmlContent, elementType, tempList)

	return listLinks
}

func (s AnalyseService) IsThereLoginForm(inputs []*html.Node, forms []*html.Node) bool {

	for _, form := range forms {
		if form.Type == html.ElementNode && form.Data == "form" {
			//this attr could be, method, id, action, class. So when we have a form related login keywords, that means we could have login form
			for _, attribute := range form.Attr {
				if strings.Contains(attribute.Val, "login") ||
					strings.Contains(attribute.Val, "signin") ||
					strings.Contains(attribute.Val, "sign-in") ||
					strings.Contains(attribute.Val, "user") {
					return true
				}
			}
		}
	}

	for _, input := range inputs {
		if input.Type == html.ElementNode && input.Data == "input" {
			// password is an important keyword here, because input type could be the password directly.
			for _, attribute := range input.Attr {
				if strings.Contains(attribute.Val, "password") ||
					attribute.Val == "pwd" {
					return true
				}
			}
		}
	}

	return false
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

func (s AnalyseService) CountOfAccessibleUrls(urls []string, baseUrl string) int {
	var wg sync.WaitGroup
	var count int
	externalIdentifiers := []string{
		"http://",
		"https://",
	}
	wg.Add(len(urls))
	c := make(chan int, len(urls))

	for _, url := range urls {
		linkForGo := url

		isExternal := false
		for _, identifier := range externalIdentifiers {
			if strings.HasPrefix(strings.TrimSpace(url), identifier) {
				isExternal = true
				break
			}
		}
		if !isExternal {
			if baseUrl[len(baseUrl)-1:] != "/" {
				baseUrl = baseUrl + "/"
			}
			linkForGo = baseUrl + url
		}

		go func() {
			defer wg.Done()

			response, err := http.Get(linkForGo)
			if err != nil {
				c <- 0
				log.Println(err)
			} else if response.StatusCode > 399 {
				c <- 0
			} else {
				c <- 1
				defer response.Body.Close()
			}
		}()
	}

	wg.Wait()
	close(c)
	for i2 := range c {
		count += i2
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
