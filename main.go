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
		org       string
		stateFile string
	)

	flag.StringVar(&org, "org", "", "Override the default organization")
	flag.StringVar(&stateFile, "tfstate", "", "Path to look for the tfstate file")
	flag.Parse()

	cmd := snyk_iac_capture.Command{
		Org:       org,
		StateFile: stateFile,
	}

	os.Exit(cmd.Run())
}
