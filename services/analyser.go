package services

type IAnalyserService interface {
	GetHtmlContentOfURL(url string) string
	DetermineHTMLVersion(html string) string
	FindPageTitle(html string) string
	FindAllUrlsInPage(html string) []string
	CountOfInternalUrls(urls []string) int
	CountOfExternalUrls(urls []string) int
	CountOfInaccessibleUrls(urls []string) int
	FindHeadingsAndCounts(html string) map[string]int
	DoesPageContainLoginForm(html string) bool
}
