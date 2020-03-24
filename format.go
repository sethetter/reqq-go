package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func FormatResponse(res *http.Response) (string, error) {
	out := res.Status + "\n"

	for k, v := range res.Header {
		out = out + fmt.Sprintf("%s: %s\n", k, v)
	}

	out = out + "\n"

	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return out, err
	}

	body, err := FormatBody(res.Header.Get("content-type"), rawBody)
	if err != nil {
		return out, err
	}
	out += body

	return out, nil
}

func FormatBody(cType string, body []byte) (string, error) {
	switch cType {
	case "application/json":
		b := &bytes.Buffer{}
		if err := json.Indent(b, body, "", "  "); err != nil {
			return "", err
		}
		return b.String(), nil
	}
	return string(body), nil
}
