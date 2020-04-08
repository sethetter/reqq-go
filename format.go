package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func FormatResponse(res *http.Response) (string, error) {
	out := res.Status + "\n"

	for k, vs := range res.Header {
		// Same header can have multiple entries.
		for _, v := range vs {
			out = out + fmt.Sprintf("%s: %s\n", k, v)
		}
	}

	out = out + "\n"

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return out, err
	}

	body := FormatBody(res.Header.Get("content-type"), rawBody)
	out += body

	return out, nil
}

func FormatBody(cType string, body []byte) string {
	if strings.Contains(cType, "application/json") {
		b := &bytes.Buffer{}
		if err := json.Indent(b, body, "", "  "); err != nil {
			return string(body)
		}
		return b.String()
	}

	return string(body)
}
