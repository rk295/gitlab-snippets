package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

const (
	defaultHost = "gitlab.com"
	apiPath     = "api/v4/snippets"

	tokenEnvVar = "GITLAB_TOKEN"
	hostEnvVar  = "GITLAB_HOST"
)

var (
	debug       bool
	description string
	file        string
	host        string
	title       string
	visibility  string
)

func init() {
	flag.StringVarP(&title, "title", "t", "snippet", "Title for the snippet")
	flag.StringVarP(&description, "description", "d", "", "Description for the snippet")
	flag.StringVarP(&visibility, "visibility", "v", "internal", "Visibility of the snippet. Possible values are: private, public, internal")
	flag.StringVarP(&file, "file", "f", "", "File to read, defaults to STDIN")
	flag.StringVarP(&host, "host", "h", defaultHost, "Host to connect to")
	flag.BoolVar(&debug, "debug", false, "sets log level to debug")
	flag.Parse()
}

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	err := checkVisibility()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	token, err := getGitlabToken()
	if err != nil {
		panic(err)
	}

	message, filename, err := getContent()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	snippet := Snippet{
		Title:       title,
		Description: description,
		Visibility:  visibility,
		Files: []File{
			{
				Content:  message,
				FilePath: filename,
			},
		},
	}

	snippetJSON, err := json.MarshalIndent(snippet, "", "  ")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, getURL(), bytes.NewBuffer(snippetJSON))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("PRIVATE-TOKEN", token)

	client := &http.Client{}
	log.Print("making request")
	resp, err := client.Do(req)
	if err != nil {
		log.Print("error making request")
		panic(err)
	}
	defer req.Body.Close()
	log.Printf("request made, response code: %v", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print("failed reading response body", err)
		}
		log.Print(string(b))
		os.Exit(1)
	}

	response := new(SnippetCreateResponse)
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Print("error decoding response JSON")
		panic(err)
	}

	fmt.Printf("%s\n", response.WebURL)
}

func getGitlabToken() (string, error) {
	t := os.Getenv(tokenEnvVar)
	if t == "" {
		return "", fmt.Errorf("couldn't read token from %s", tokenEnvVar)
	}
	return t, nil
}

func getContent() (string, string, error) {

	if file != "" {
		log.Printf("attempting to read file %s", file)
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return "", "", err
		}
		return string(f), file, nil
	}

	log.Print("no file given, trying stdin")
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", "", err
	}

	if fi.Size() > 0 {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", "", err
		}
		return string(bytes), "-", nil
	}

	log.Print("no file given and nothing on stdin")
	return "", "", errors.New("must pass either --file or stdin")
}

func checkVisibility() error {
	switch visibility {
	case "internal", "private", "public":
		return nil
	}
	return fmt.Errorf("%s is not a possible visibility, options are public, private, internal", visibility)
}

func getURL() string {
	// Allows the env var to take precedence over the command line flag
	hostname := os.Getenv(hostEnvVar)
	if hostname == "" {
		log.Printf("%s unset, defaulting to %s", hostEnvVar, host)
		hostname = host
	}
	url := fmt.Sprintf("https://%s", path.Join(hostname, apiPath))
	log.Printf("will use url: %s", url)
	return url
}
