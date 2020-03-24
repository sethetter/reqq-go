# reqq

CLI for sending HTTP requests defined in text files.

## Installing

```
$ go install github.com/sethetter/reqq
```

## Usage

Without a subcommand, `reqq`'s default behavior is to issue the defined http request.

```
$ reqq -e reqq/envs/test.json reqq/plan-search/nevada-single-family.txt
```

This will..
- Get the request file at the specified path
- Look up the `reqq/envs/test.json` file the variables to inject into the http request template.
- Issue the request and output the response data.

## Request File Format

```
POST http://the.url.to/send/to?query=here
x-header-a: header values here
x-header-b: yet-another-header
body
content
follows
to EOF
```

# Development

This project is under active development!

## TODO

- Environment files.
- Better output.
  - More content-type formats?
  - Colors?
- Comments in request files.
