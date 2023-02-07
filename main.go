package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	capture "github.com/snyk/snyk-iac-capture/cmd/snyk-iac-capture"
)

var (
	version string
	commit  string
)

func main() {
	flag.BoolP("debug", "d", false, "Show debug information")
	//flag.Bool("http-tls-skip-verify", false, "If set, skip client validation of TLS certificates") Disabled for the closed beta
	flag.String("api-rest-url", "https://api.snyk.io", "Url for Snyk REST API")
	flag.String("path", ".", "Path to look for Terraform state files (can be a file, a directory or a glob pattern)")
	flag.String("api-rest-token", "", "Auth token for the API Usage (Required)")
	flag.String("org", "", "Organization public id (Required)")
	flag.Parse()

	// normalize flag with - in the name to make it easier to match with env
	f := flag.CommandLine
	normalizeFunc := f.GetNormalizeFunc()
	f.SetNormalizeFunc(func(fs *flag.FlagSet, name string) pflag.NormalizedName {
		result := normalizeFunc(fs, name)
		name = strings.ReplaceAll(string(result), "-", "_")
		return pflag.NormalizedName(name)
	})

	viper.SetEnvPrefix("SNYK_IAC_CAPTURE")
	must(viper.BindPFlags(flag.CommandLine))
	must(viper.BindEnv("debug"))
	must(viper.BindEnv("org"))
	must(viper.BindEnv("path"))
	must(viper.BindEnv("http_tls_skip_verify"))
	must(viper.BindEnv("api_rest_url"))
	must(viper.BindEnv("api_rest_token"))

	logrus.SetLevel(logrus.WarnLevel)
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Debugf("snyk-iac-capture %s (%s)", version, commit)

	command := capture.Command{
		Org:               viper.GetString("org"),
		StatePath:         viper.GetString("path"),
		HTTPTLSSkipVerify: viper.GetBool("http_tls_skip_verify"),
		APIURL:            viper.GetString("api_rest_url"),
		APIToken:          viper.GetString("api_rest_token"),
		ExtraSSlCerts:     os.Getenv("NODE_EXTRA_CA_CERTS"), // we still want to read this one without prefix to work with snyk node cli integration
	}

	os.Exit(command.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
