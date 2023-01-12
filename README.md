# IaC Capture

This repository contains the command to upload filtered states artifacts to Snyk's api. 

## Configuration
Update your Git configuration so the Go toolchain will use Git over SSH instead
of HTTPS:

```
[url "ssh://git@github.com/"]
	insteadOf = https://github.com/
```

Add to your `PATH` the directory where Go install executables by default:

```sh
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Build

Build and install the project:

```
go install
```

The executable `snyk-iac-capture` should now be in your `PATH`.

### Run

## Build

We are using `goreleaser` to build and release binaries. We have wrapped up
everything in a Makefile for easier development process.

### Development binaries

To build a binary, run:

```shell
make build
```

This command will generate a `dist` directory in the root of this repo, which
will contain a binary with the current's machine architecture.

### Release binaries

To build and release binaries, run:

```shell
make release
```

This command will generate a `dist` directory in the root of this repo, which
will contain binaries for different architectures.

## Release

The release of a new version is gated behind a CircleCI approval mechanism, this
way we can select when we would like to release a new version of the
`snyk-iac-capture` binary.

To approve the release of a new version:

1. Go to <TODO TBD>
2. Open the workflow of the latest merge
3. Click on `Approve Release` and hit the `Approve` button

To reject the release of a new version:

1. Go to <TODO TBD>
2. Cancel the workflow of the latest merge

You can check if the binaries have been uploaded to the S3 bucket like so:

```
aws-vault exec <user> -- aws s3 ls s3://snyk-assets/cli/iac/capture/v*.*.*/
```

[snyk-ci]: https://app.circleci.com/pipelines/github/snyk/cli?branch=master
[snyk-latest-release]: https://github.com/snyk/cli/releases/latest
[aws-vault]: https://github.com/99designs/aws-vault
