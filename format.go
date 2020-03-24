package main

import (
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return out, err
	}

	out += string(body)

	return out, nil
}

var example string = `
200 OK
Header-Val: yonder, donder, shmonder

body content!
`
