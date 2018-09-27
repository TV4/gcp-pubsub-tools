# gcp-pubsub-subscribe

[![Build Status](https://travis-ci.org/TV4/gcp-pubsub-tools.svg?branch=master)](https://travis-ci.org/TV4/gcp-pubsub-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/gcp-pubsub-tools)](https://goreportcard.com/report/github.com/TV4/gcp-pubsub-tools)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/gcp-pubsub-tools#license)

`gcp-pubsub-subscribe` reads from the specified Pub/Sub subscription and writes each
message followed by a newline to `stdout`.

## Installation
```
go get -u github.com/TV4/pubsub-tools/cmd/gcp-pubsub-subscribe
```

## Usage
```
gcp-pubsub-subscribe -credentials=<...> -project=<...> -subscription=<...>

  -credentials string
    	path to a GCP credentials file
  -project string
    	Pub/Sub project ID
  -quiet
    	do not print messages (side-effect: more speed)
  -subscription string
    	Pub/Sub subscription name

GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials
```

## License
Copyright (c) 2017-2018 Bonnier Broadcasting / TV4

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
