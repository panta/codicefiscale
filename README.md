# panta/codicefiscale

Small library to check and decode Italian "Codice Fiscale" from [Go](https://golang.org).

## Installation

First, use `go get` to install the latest version of the library:

```bash
go get -u github.com:panta/codicefiscale
```

Then include the library in your code:

```go
import "panta/codicefiscale"
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/panta/codicefiscale"
)

func main() {
	cf, err := codicefiscale.Decode("RSSMRA77L18H501W")
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Birth date: %s place: %s (%s)",
		cf.BirthDate.Format(time.RFC3339),
		cf.BirthPlaceName, cf.BirthPlace.Provincia.Nome)
}
```
