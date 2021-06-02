# gcp-gcs

[![Build Status](https://travis-ci.com/TV4/gcp-tools.svg?branch=master)](https://travis-ci.com/TV4/gcp-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/gcp-tools)](https://goreportcard.com/report/github.com/TV4/gcp-tools)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/gcp-tools#license)

`gcp-gcp` reads and writes from/to
[Google Cloud Storage](https://cloud.google.com/storage/).

## Installation
```
go install github.com/TV4/gcp-tools/cmd/gcp-gcs@latest
```

## Usage
```
gcp-gcs [-credentialsfile=<...>|-credentialsjson=<...>] -bucket=<...> <command>

  -bucket string
        bucket name
  -credentialsfile string
        path to a GCP credentials file
  -credentialsjson string
        json string of a GCP credentials file content

  Commands:
    ls       [<prefix>]              Lists bucket objects, optionally filtered by the given prefix
    download <object> [<object>...]  Downloads given object(s) from the bucket
    upload   <file> [<file>...]      Uploads given file(s) to the bucket
    read     <object>                Reads a bucket object to stdout
    write    <object>                Writes a bucket object from stdin
    rm       <object> [<object>...]  Deletes the given object(s) from the bucket


GCP credentials file:
  https://developers.google.com/identity/protocols/application-default-credentials
```

## License
Copyright (c) 2019 TV4

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
