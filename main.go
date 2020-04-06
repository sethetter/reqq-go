package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

var usage string = `reqq: help

reqq path/to/request/file.txt
`

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Path to the JSON env file providing the data to fill in to request template vars.",
			},
		},
		Action: mainAction,
	}

	app.Run(os.Args)
}

func mainAction(c *cli.Context) error {
	reqPath := c.Args().Get(0)
	envPath := c.String("env")

	reqF, envF, err := getReqAndEnvFiles(reqPath, envPath)
	if err != nil {
		fmt.Printf("failed getting req and env files: %v", err)
		return err
	}

	req, err := NewRequest(reqF, envF)
	if err != nil {
		fmt.Printf("failed to parse request: %v", err)
		return err
	}

	res, err := req.Send(http.DefaultClient)
	if err != nil {
		fmt.Printf("failed to execute http request: %v", err)
		return err
	}

	out, err := FormatResponse(res)
	if err != nil {
		fmt.Printf("failed formatting response: %v", err)
		return err
	}

	fmt.Printf("%s\n", out)

	return nil
}

func getReqAndEnvFiles(reqPath string, envPath string) (io.Reader, io.Reader, error) {
	reqF, err := os.Open(reqPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.New("request file does not exist")
		}
		return nil, nil, err
	}

	var envF io.Reader
	if envPath != "" {
		f, err := os.Open(envPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil, errors.New("environment file does not exist")
			}
			return nil, nil, err
		}
		envF = f
	}

	return reqF, envF, nil
}
