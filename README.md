# reqq

CLI for sending HTTP requests defined through config files defined in a local project folder.

- `.reqq/envs/` defines available environments for the project.
- `.reqq/reqs/` defines the requests that can be performed.

## Commands

Without a subcommand, `reqq`'s default behavior is to issue the defined http request.

```
$ reqq -e test plan-search/nevada-single-family.txt
```

This will..
- Get the http request definition from `.reqq/reqs/plan-search/nevaga-single-family.json`.
- Look up the `.reqq/envs/test.json` file the variables to inject into the http request template.
- Issue the request and output the response data.

### Subcommands

- `list` will show a list of defined requests.
- `envs` will list the defined environments.
-`edit <request>` will run `$EDITOR .req/reqs/<request>.txt`
-`edit-env <env>` will run `$EDITOR .req/envs/<env>.json`

## Request File Format

```
http://the.url.to/send/to?query=here
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

- [ ] (WIP) Request file parsing
- [ ] Argument parsing
- [ ] Load request from argument
- [ ] Request execution
- [ ] CLI `list` command
- [ ] CLI `edit` command
- [ ] Environment support
- [ ] CLI `edit-env` command
