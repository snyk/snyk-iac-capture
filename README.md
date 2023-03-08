<p align="center">
  <img src="https://snyk.io/style/asset/logo/snyk-print.svg" />
</p>

# Snyk IaC Capture CLI Extension

## Overview

This module implements the Snyk CLI Extension to capture a Terraform state, filter it and send it to Snyk API.


## Filtering

We use an allowlist to ensure we only send expected and non-sensitive information to the Snyk API.

You can find the expected fields in [state.go](internal/terraform/state.go) and [statefilter.go](internal/filtering/statefilter.go).

## Usage

This repository produces a standalone binary for debugging purposes. We advise to use this command as part of [Snyk CLI](https://github.com/snyk/cli).

This is the usage for the standalone binary:
```
Usage of snyk-iac-capture:
      --api-rest-token string   Auth token for the API Usage (Required)
      --api-rest-url string     Url for Snyk REST API (default "https://api.snyk.io")
  -d, --debug                   Show debug information (Optional default false)
      --org string              Organization public id (Required)
      --path string             Path to look for Terraform state files (can be a file, a directory or a glob pattern) (Optional default ".")
      --stdin                   Read Terraform state from the standard input instead of path (Optional)
```

