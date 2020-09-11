# Probe
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/gomicro/probe/Build/master)](https://github.com/gomicro/probe/actions?query=workflow%3ABuild+branch%3Amaster)
[![Go Reportcard](https://goreportcard.com/badge/github.com/gomicro/probe)](https://goreportcard.com/report/github.com/gomicro/probe)
[![GoDoc](https://godoc.org/github.com/gomicro/probe?status.svg)](https://godoc.org/github.com/gomicro/probe)
[![License](https://img.shields.io/github/license/gomicro/probe.svg)](https://github.com/gomicro/probe/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/gomicro/probe.svg)](https://github.com/gomicro/probe/releases/latest)

Probe is for use inside of Docker Containers build on the `scratch` image.  Given the lack of any tools to do something as simple as a curl command, it is not possible to add a `HEALTHCHECK` to your Dockerfile.  Probe provides an ultra stripped down curl command to bundle with your super slim containers.

# Installation

## Precompiled Binary

See the [Latest Release](https://github.com/gomicro/probe/releases/latest) page for a download link to the binary compiled for your system.

## From Source

Requires Golang version 1.14 or higher

```
go get -u github.com/gomicro/probe
go install github.com/gomicro/probe
```

# Usage

Probe has one option.  You provide it with a URL to ping.  If it is not able to get a HTTP OK back from the URL, it will return an exit code.
```
$ probe http://localhost:4567/v1/status
```

# Versioning
The cli will be versioned in accordance with [Semver 2.0.0](http://semver.org).  See the [releases](https://github.com/gomicro/probe/releases) section for the latest version.  Until version 1.0.0 the cli is considered to be unstable.

# License
See [LICENSE.md](./LICENSE.md) for more information.
