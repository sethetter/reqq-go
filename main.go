package main

import (
	"errors"
	"fmt"
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
		Action: func(c *cli.Context) error {
			f, err := os.Open(c.Args().Get(0))
			if err != nil {
				if os.IsNotExist(err) {
					return errors.New("request file does not exist")
				}
				return err
			}

			req, err := NewRequest(f)
			if err != nil {
				return err
			}

			// send it
			res, err := req.Send(http.DefaultClient)
			if err != nil {
				return err
			}

			// TODO: transform the output
			out, err := FormatResponse(res)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", out)

			return nil
		},
	}

	app.Run(os.Args)
}
