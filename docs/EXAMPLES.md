# Examples

## vcs

Install the vcs package: `go get github.com/bit8bytes/toolbox/vcs` and use it:

```go
package main

import (
	"fmt"

	"github.com/bit8bytes/toolbox/vcs"
)

func main() {
	fmt.Println(vcs.Version())
}
```

Build the binary using `go build -o=./bin/vcs .` and run `./bin/vcs`

If you run `go run .` the output will be just `-`. The version is only available in the build binary and needs a version control system like Git already setup.

## Middleware

### GZIP

### CORS

## Responder

### JSON

## validator

Install the validator package: `go get github.com/bit8bytes/toolbox/validator` and use it:

```go
package main

import (
	"fmt"

	"github.com/bit8bytes/toolbox/validator"
)

func main() {
	user := "bit9bytes"

	v := validator.New()

	// Check multiple conditions to be true.
	// The len of user is greater then 0, therefore not error.
	v.Check(len(user) != 0, "name", "Name cannot be empty")
	// The user isn't bit8bytes, therefore error.
	v.Check(user == "bit8bytes", "name", "Name must be 'bit8bytes'")

	// If any of this checks is not valid, the v.Valid() will return falls.
	// Therefore we map over the v.Errors.
	if !v.Valid() {
		for field, msg := range v.Errors {
			// The error will be: "Name must be 'bit8bytes'"
			fmt.Printf("key: %s: msg: %s\n", field, msg)
		}
		return
	}

	fmt.Println("Validation passed!")
}
```

Run `go run .`, the output will be:

```bash
$ go run .
key: name: msg: Name must be 'bit8bytes'
```
