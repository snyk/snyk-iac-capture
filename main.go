package main

import (
	"flag"
	"log"
	"os"

	snyk_iac_capture "github.com/snyk/snyk-iac-capture/cmd/snyk-iac-capture"
)

var (
	version string
	commit  string
)

func main() {
	if version != "" && commit != "" {
		log.Printf("snyk-iac-capture %s (%s)", version, commit)
	}

	var (
		org               string
		stateFile         string
		httpTLSSkipVerify bool
	)

	flag.StringVar(&org, "org", "", "Override the default organization")
	flag.StringVar(&stateFile, "tfstate", "", "Path to look for the tfstate file")
	flag.BoolVar(&httpTLSSkipVerify, "http-tls-skip-verify", false, "If set, skip client validation of TLS certificates")
	flag.Parse()

	apiURL := os.Getenv("SNYK_IAC_CAPTURE_API_REST_URL")
	apiToken := os.Getenv("SNYK_IAC_CAPTURE_API_REST_TOKEN")

	cmd := snyk_iac_capture.Command{
		Org:               org,
		StateFile:         stateFile,
		HTTPTLSSkipVerify: httpTLSSkipVerify,
		APIURL:            apiURL,
		APIToken:          apiToken,
	}

	os.Exit(cmd.Run())
}
