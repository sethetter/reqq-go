package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	is "github.com/go-ozzo/ozzo-validation/is"
)

// Request represents an HTTP request to be sent.
type Request struct {
	// TODO: this should just wrap the net/http.Request struct
	Method  string
	URL     string
	Headers http.Header
	Body    string
}

type ParseError struct {
	Msg  string
	Err  error
	Part string
	Hint string
}

func (e *ParseError) Error() string {
	errStr := fmt.Sprintf("parse error: %s\n", e.Msg)

	if e.Part != "" {
		errStr = errStr + fmt.Sprintf("\nFailed on: %s\n", e.Part)
	}

	if e.Hint != "" {
		errStr = errStr + fmt.Sprintf("\nHinnt: %s\n", e.Hint)
	}

	return errStr
}

func (e *ParseError) Unwrap() error { return e.Err }

func NewRequest(reqR io.Reader, envR io.Reader) (Request, error) {
	if envR != nil {
		return parseReqWithEnv(reqR, envR)
	}
	return parseReq(reqR)
}

func parseReqWithEnv(reqR io.Reader, envR io.Reader) (Request, error) {
	env, err := ioutil.ReadAll(envR)
	if err != nil {
		return Request{}, &ParseError{Msg: "failed to read env file", Err: err}
	}

	req, err := ioutil.ReadAll(reqR)
	if err != nil {
		return Request{}, &ParseError{Msg: "failed to read req file", Err: err}
	}

	// parse the env file as JSON
	envMap := map[string]string{}
	if err := json.Unmarshal(env, &envMap); err != nil {
		return Request{}, &ParseError{
			Msg:  "failed parsing env file as JSON",
			Err:  err,
			Hint: "env files should only have string values!",
		}
	}

	// parse the request file as a template
	reqT := template.New("request")
	if reqT, err = reqT.Parse(string(req)); err != nil {
		return Request{}, &ParseError{
			Msg: "failed to parse request file as a template",
			Err: err,
		}
	}

	// execute the template with the env data
	var parsedReq bytes.Buffer
	if err := reqT.Execute(&parsedReq, envMap); err != nil {
		return Request{}, &ParseError{
			Msg: "failed compiling request with env data",
			Err: err,
		}
	}

	return parseReq(&parsedReq)
}

func parseReq(reqR io.Reader) (Request, error) {
	var req Request

	lines := bufio.NewScanner(reqR)

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
		req.Headers.Add(parts[0], strings.TrimSpace(parts[1]))
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
	req, err := http.NewRequest(r.Method, r.URL, bodyR)
	if err != nil {
		return req, err
	}
	req.Header = r.Headers
	return req, nil
}

func (r *Request) Send(c *http.Client) (*http.Response, error) {
	req, err := r.Build()
	if err != nil {
		return &http.Response{}, err
	}
	return c.Do(req)
}
