package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

// Request represents an HTTP request to be sent.
type Request struct {
	// TODO: should this just wrap the net/http.Request struct?
	Method  string
	URL     string
	Headers http.Header
	Body    string
}

type ParseError struct {
	Msg  string
	Err  error
	Part string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error: %s\n\nFailed on: %s\n", e.Msg, e.Part)
}

func (e *ParseError) Unwrap() error { return e.Err }

func NewRequest(r io.Reader) (Request, error) {
	var req Request

	lines := bufio.NewScanner(r)

	// Parse method and URL from first line.
	lines.Scan()
	if err := lines.Err(); err != nil {
		return req, &ParseError{Msg: "failed to read first line", Err: err}
	}

	first := lines.Text()
	parts := strings.SplitN(strings.Trim(first, "\n"), " ", 2)
	if len(parts) != 2 {
		return req, &ParseError{Msg: "first line must include method and url", Part: first}
	}

	req.Method = strings.ToUpper(parts[0])
	if err := validateMethod(req.Method); err != nil {
		return req, &ParseError{Msg: "invalid method", Err: err, Part: req.Method}
	}

	req.URL = parts[1]
	if err := validateURL(req.URL); err != nil {
		return req, &ParseError{Msg: "invalid url", Err: err, Part: req.URL}
	}

	// Parse header lines.
	headerRE := regexp.MustCompile(`^[\w-]+: .*`)
	req.Headers = make(http.Header)
	for lines.Scan() {
		line := lines.Text()

		if !headerRE.MatchString(line) {
			// This should be first line of body.
			req.Body = line
			break
		}

		parts = strings.SplitN(line, ": ", 2)
		req.Headers.Add(parts[0], strings.Trim(parts[1], "\n"))
	}
	if err := lines.Err(); err != nil {
		return req, &ParseError{Msg: "failed reading header lines", Err: err}
	}

	// Remaining lines are the rest of the body.
	for lines.Scan() {
		req.Body += "\n" + lines.Text()
	}
	if err := lines.Err(); err != nil {
		return req, &ParseError{Msg: "failed reading body lines", Err: err}
	}

	return req, nil
}

func validateMethod(method string) error {
	return validation.Validate(
		method,
		validation.Required,
		validation.In("OPTIONS", "HEAD", "GET", "POST", "PUT", "DELETE", "PATCH"),
	)
}

func validateURL(url string) error {
	return validation.Validate(
		url,
		validation.Required,
		is.URL,
	)
}

func (r *Request) Build() (*http.Request, error) {
	var bodyR io.Reader
	if r.Body != "" {
		bodyR = strings.NewReader(r.Body)
	}
	return http.NewRequest(r.Method, r.URL, bodyR)
}

func (r *Request) Send(c *http.Client) (*http.Response, error) {
	req, err := r.Build()
	if err != nil {
		return &http.Response{}, err
	}
	return c.Do(req)
}
