# Examples

## vcs

Install the vcs package: `go get github.com/bit8bytes/toolbox/vcs` and use it:

```go
package main

import (
    "fmt"
    "github.com/bit8bytes/toolbox/vcs"
)

fun main() {
    fmt.Println(vcs.Version())
}
```

Build the binary: `go build -o=./bin/vcs ./cmd` and run `./bin/vcs`

If you run `go run .` the output will be just `-`. The version is only available in the build binary.

## Middleware

### GZIP

### CORS

## Responder

### JSON

## Validator
