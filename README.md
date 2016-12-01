# sentrytool

sentrytool is a Go library and command-line interface to [Apache Sentry](http://sentry.apache.org/).

The tool and library can be used to interface with non-kerberized Sentry daemon.

## Motivation

- Provide command-line access to Sentry Thrift API
- Provide simple command-line tool that doesn't require any infrastructure
  to run

## Features

- Flexible configuration
  * Environment variables
  * Config files
  * Command-line arguments
- Ease of use
  * Intuitive interface
  * Flexible way of specifying parameters
- Easy deployment
  * Single binary
  * Easy to build and install (using `go get`)
  
## Installation

Standard `go get`:

```
$ go get github.com/akolb1/sentrytool
```

## Usage & Example

* [sentrytool](doc/sentrytool.md) - Command-line interface to Apache Sentry
* [API](sentryapi/README.md) - API Usage and examples
