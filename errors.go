package cronofy

import (
	"fmt"
	"net/http"
	"net/url"
)

// responseError provides a structured type for errors caused by
// external requests.
type responseError struct {
	// statusCode is the status code of the response.
	statusCode int

	// url is the requested resource.
	url *url.URL
}

// HTTPStatusCode returns the status code of the response.
func (e *responseError) HTTPStatusCode() int { return e.statusCode }

// URL returns the request URL.
func (e *responseError) URL() *url.URL { return e.url }

// Error satisfies the error interface returning the underlying error.
func (e *responseError) Error() string {
	return fmt.Sprintf("restclient: url=%s returned status=%d %s", e.url, e.statusCode, http.StatusText(e.statusCode))
}
