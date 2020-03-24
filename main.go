package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var usage string = `reqq: help

reqq path/to/request/file.txt
`

func main() {
	// parse the first arg, check for existence of file at .reqq/...
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		panic(usage)
	}

	// read the file, parse it into a request
	f, err := os.Open(args[0])
	if err != nil {
		if os.IsNotExist(err) {
			panic("request file does not exist")
		}
		panic(err)
	}

	req, err := NewRequest(f)
	if err != nil {
		panic(err)
	}

	// send it
	res, err := req.Send(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	// TODO: transform the output
	out, err := FormatResponse(res)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", out)
}
