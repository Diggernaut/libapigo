# libapigo

Diggernaut Official API Library for Go
This repo still in development.

#### install:

go get github.com/Diggernaut/libapigo

#### usage:
```go
package main

import (
	"fmt"

	diggernaut "github.com/Diggernaut/libapigo"
)

func main() {
	dig := diggernaut.New("f98c1dc37033a8b1755f839685062cf422221111")
	projects := dig.GetProjects()
	fmt.Println(projects)
	//prints [{1 ProjectName  []}]
}
```
License MIT.
