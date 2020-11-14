# beertime

beer time is a binary time format that indicates if it's beer time or if you should get back to work!

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/nuttmeister/beertime"
)

func main() {
	// Get current local time.
	now := time.Now().Local()

	// Get beer time.
	beer := beertime.Now(now)
	dur := beertime.Duration(now)

	fmt.Printf("current beer time: %t. time to next: %s\n", beer, dur)
}

```