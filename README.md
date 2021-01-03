# hex

[![GoDoc](https://godoc.org/github.com/VKoptev/hex?status.svg)](https://godoc.org/github.com/VKoptev/hex)

Realization of hexagonal grids

This library is representation of hexagonal grids that describe in [article from Red Blob Games]

[article from Red Blob Games]: https://www.redblobgames.com/grids/hexagons/

# How to use

```go
package main

import (
	"fmt"
	
	"github.com/VKoptev/hex"
)

func main() {
	h := hex.New(1, 2)

	fmt.Printf("equal: %v", h.Equal(hex.ZE))
	fmt.Printf("to east: %v", h.Add(hex.EE))
	fmt.Printf("to west: %v", h.Sub(hex.EE))
	fmt.Printf("no way: %v", hex.ZE.Mul(10))
	// ...
}
```
