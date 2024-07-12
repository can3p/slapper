# slapper

__Simple load testing tool with real-time updated histogram of request timings__

![slapper](https://raw.githubusercontent.com/ikruglov/slapper/master/img/example.gif)

## Interface

![interface](https://raw.githubusercontent.com/ikruglov/slapper/master/img/interface.png)

## Usage
```bash
$ ./slapper run --help
Usage of ./slapper:
  -H value
    	HTTP header 'key: value' set on all requests. Repeat for more than one header.
  --base64body
    	Bodies in targets file are base64-encoded
  --maxY duration
    	max on Y axe (default 100ms)
  --minY duration
    	min on Y axe (default 0ms)
  --rate uint
    	Requests per second (default 50)
  --targets string
    	Targets file
  --timeout duration
    	Requests timeout (default 30s)
  --workers uint
    	Number of workers (default 8)

```

## Key bindings
* q, ctrl-c - quit
* r - reset stats
* k - increase rate by 100 RPS
* j - decrease rate by 100 RPS

## Targets syntax

The targets file is line-based. Its syntax is:

	HTTP_METHOD url
	H Header: Value
	$ body

Header lines are optional.

The body line is optional. The rules for what is considered to be a body
line are:

1. If something starts with `$ ` (dollar-sign and space), it's a body
2. If the line is literally `{}`, it's an empty body

A missing body line is taken to mean an empty request body. Point (2) is there
for backwards-compatibility.

## Installation

### Install Script

Download `slapper` and install into a local bin directory.

#### MacOS, Linux, WSL

Latest version:

```bash
curl -L https://raw.githubusercontent.com/can3p/slapper/master/generated/install.sh | sh
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/can3p/slapper/master/generated/install.sh | sh -s 0.0.4
```

The script will install the binary into `$HOME/bin` folder by default, you can override this by setting
`$CUSTOM_INSTALL` environment variable

### Manual download

Get the archive that fits your system from the [Releases](https://github.com/can3p/slapper/releases) page and
extract the binary into a folder that is mentioned in your `$PATH` variable.

## Notes

The project has been scaffolded with the help of [kleiner](https://github.com/can3p/kleiner)

## Acknowledgement
* Idea and initial implementation is by @sparky
* This module was originally developed for Booking.com.
