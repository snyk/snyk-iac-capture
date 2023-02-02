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
	flag.Bool("debug", false, "Override the default organization")
	flag.String("org", "", "Organization public id")
	flag.String("path", "", "Path to look for the terraform state files (can be a file, a directory or a glob pattern)")
	flag.Bool("http-tls-skip-verify", false, "If set, skip client validation of TLS certificates")
	flag.String("api-rest-url", "", "Path to look for the tfstate file")
	flag.String("api-rest-token", "api-rest-token", "Path to look for the tfstate file")
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
	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("debug"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("org"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("path"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("http_tls_skip_verify"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("api_rest_url"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("api_rest_token"); err != nil {
		panic(err)
	}

	logrus.SetLevel(logrus.WarnLevel)
	if (version != "" && commit != "") || viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("snyk-iac-capture %s (%s)", version, commit)
	}
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
