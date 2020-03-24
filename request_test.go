package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		desc   string
		in     string
		err    error
		expect Request
	}{
		{
			desc: "happy path",
			// TODO: switch to fixture files?
			in: strings.Trim(`
POST https://some.url/yaknow?sup=yup
x-header-thing: some header value
x-header-2: another header
this is the body content
yadayadayada yadayadayada`, "\n"),
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
			err:    errors.New("invalid method"),
			expect: Request{},
		},
		{
			desc:   "invalid URL",
			in:     "POST not_a*valid!url",
			err:    errors.New("invalid url"),
			expect: Request{},
		},
		{
			desc:   "invalid first line",
			in:     "https://just.a.url.com/lol",
			err:    errors.New("first line must include"),
			expect: Request{},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d_%s", i, tc.desc), func(t *testing.T) {
			req, err := NewRequest(strings.NewReader(tc.in))

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
