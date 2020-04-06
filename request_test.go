package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		desc   string
		in     string
		env    string
		err    error
		expect Request
	}{
		{
			desc: "happy path, no env",
			// TODO: switch to fixture files?
			in: strings.Trim(`
POST https://some.url/yaknow?sup=yup
x-header-thing: some header value
x-header-2: another header
this is the body content
yadayadayada yadayadayada`, "\n"),
			env: "",
			err: nil,
			expect: Request{
				Method: "POST",
				URL:    "https://some.url/yaknow?sup=yup",
				Headers: map[string][]string{
					"X-Header-Thing": []string{"some header value"},
					"X-Header-2":     []string{"another header"},
				},
				Body: "this is the body content\nyadayadayada yadayadayada",
			},
		},
		{
			desc:   "invalid method",
			in:     "SWHAT https://valid.url.com/freal",
			env:    "",
			err:    errors.New("invalid method"),
			expect: Request{},
		},
		{
			desc:   "invalid URL",
			env:    "",
			in:     "POST not_a*valid!url",
			err:    errors.New("invalid url"),
			expect: Request{},
		},
		{
			desc:   "invalid first line",
			env:    "",
			in:     "https://just.a.url.com/lol",
			err:    errors.New("first line must include"),
			expect: Request{},
		},
		{
			desc: "happy path, with env",
			env:  `{"baseUrl": "https://just.a.url.com"}`,
			in:   "GET {{ .baseUrl }}/lol",
			err:  nil,
			expect: Request{
				Method:  "GET",
				URL:     "https://just.a.url.com/lol",
				Headers: http.Header{},
				Body:    "",
			},
		},
		{
			desc:   "failure: env with non-string keys",
			env:    `{"baseUrl": "https://just.a.url.com", "sup": true}`,
			in:     "GET {{ .baseUrl }}/lol",
			err:    errors.New("should only have string values"),
			expect: Request{},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tc.desc), func(t *testing.T) {
			reqR := strings.NewReader(tc.in)

			var envR io.Reader
			if tc.env != "" {
				envR = strings.NewReader(tc.env)
			}

			req, err := NewRequest(reqR, envR)

			if tc.err != nil {
				assert.Contains(t, err.Error(), tc.err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expect.Method, req.Method)
				assert.Equal(t, tc.expect.URL, req.URL)
				assert.Equal(t, tc.expect.Headers, req.Headers)
				assert.Equal(t, tc.expect.Body, req.Body)
			}
		})
	}
}
