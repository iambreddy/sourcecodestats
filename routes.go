package main

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
)

// GitHubURLPattern - Regex Pattern for GitHub request URLs
var GitHubURLPattern = `^/github/([^/]+)/([^/]+)/?$`

// GitLabURLPattern - Regex Pattern for GitLab request URLs
var GitLabURLPattern = `^/gitlab/([^/]+)/([^/]+)/?$`

// BitBucketURLPattern - Regex Pattern for BitBucket request URLs
var BitBucketURLPattern = `^/bitbucket/([^/]+)/([^/]+)/?$`

// Route - Structure to store details of a route such as URL pattern
// and it's handler
type Route struct {
	urlPattern *regexp.Regexp
	handler    http.Handler
}

// HTTPRouteHandleManager - Manages all the routes and their handlers
type HTTPRouteHandleManager struct {
	routes []*Route
}

// UniversalHandler - A single request handler which will process all the requests for the application
func (handleManager *HTTPRouteHandleManager) UniversalHandler(responseWriter http.ResponseWriter, request *http.Request) {

	// Make sure that the writer supports flushing
	if _, ok := responseWriter.(http.Flusher); !ok {
		http.Error(responseWriter, "Streaming Unsupported!", http.StatusInternalServerError)
	}

	// Stores the final handler for the request
	var handler http.Handler

	for _, route := range handleManager.routes {
		if ok := route.urlPattern != nil && route.handler != nil; ok && route.urlPattern.MatchString(request.URL.Path) {
			handler = route.handler
			break
		}
	}

	if handler == nil {
		http.NotFound(responseWriter, request)
	}

	// Set all the required HTTP headers for handling Server Side Events
	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")

	handler.ServeHTTP(responseWriter, request)
}

// InitializeRouteHandles - Initializes the handles for all available routes
func (handleManager *HTTPRouteHandleManager) InitializeRouteHandles() func(responseWriter http.ResponseWriter, request *http.Request) {

	handleManager.routes = append(handleManager.routes, GetGitHubHandler())
	handleManager.routes = append(handleManager.routes, GetGitLabHandler())
	handleManager.routes = append(handleManager.routes, GetBitBucketHandler())

	return handleManager.UniversalHandler
}

// InitializeRouteHandleManager - Initializes the HttpRouteHandleManager object
func InitializeRouteHandleManager() *HTTPRouteHandleManager {
	var routeHandler *HTTPRouteHandleManager = new(HTTPRouteHandleManager)
	routeHandler.routes = make([]*Route, 0)
	return routeHandler
}

// GithubHandler - Request handler type for GitHub
type GithubHandler struct{}

func (handler GithubHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hello from GitHub! %q", html.EscapeString(request.URL.Path))
}

// GitlabHandler - Request handler type for GitLab
type GitlabHandler struct{}

func (handler GitlabHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hello from GitLab! %q", html.EscapeString(request.URL.Path))
}

// BitbucketHandler - Request handler for BitBucket
type BitbucketHandler struct{}

func (handler BitbucketHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(responseWriter, "Hello from BitBucket! %q", html.EscapeString(request.URL.Path))
}

// GetGitHubHandler - Gets the route handler for GitHub
func GetGitHubHandler() *Route {
	var githubRegex *regexp.Regexp = regexp.MustCompile(GitHubURLPattern)
	return &Route{githubRegex, GithubHandler{}}
}

// GetGitLabHandler - Gets the route handler type for GitLab
func GetGitLabHandler() *Route {
	var gitlabRegex *regexp.Regexp = regexp.MustCompile(GitLabURLPattern)
	return &Route{gitlabRegex, GitlabHandler{}}
}

// GetBitBucketHandler - Gets the route handler for BitBucket
func GetBitBucketHandler() *Route {
	var bitbucketRegex *regexp.Regexp = regexp.MustCompile(BitBucketURLPattern)
	return &Route{bitbucketRegex, BitbucketHandler{}}
}