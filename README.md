# gitlab-snippets

Quick tool to grab a file or stdin and make a Gitlab snippet out of it. 

You must set the `GITLAB_TOKEN` variable to hold your Gitlab PAT.

## Installation

Install with the following:

```
go install github.com/rk295/gitlab-snippets@master
```

## Options

You can set `GITLAB_HOST` to override the default host (gitlab.com) and save you from passing `-f` every time.

```
snippet --help
Usage of snippet:
      --debug                sets log level to debug
  -d, --description string   Description for the snippet
  -f, --file string          File to read, defaults to STDIN
  -h, --host string          Host to connect to (default "gitlab.com")
  -t, --title string         Title for the snippet (default "snippet")
  -v, --visibility string    Visibility of the snippet. Possible values are: private, public, internal (default "internal")
```