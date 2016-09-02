# libapigo

Diggernaut Official API Library for Go

#### install:

go get github.com/Diggernaut/libapigo

#### usage:
```go
package main

import (
	"fmt"

	dig "github.com/Diggernaut/libapigo"
)

func main() {
	dig.SetUpAPIKey("f98c1dc37033a8b1755f839685062cf422221111")
	API := dig.API{}
	API.GetProjects()
	fmt.Println(API.Projects)
	//prints [{1 ProjectName  []}]
}
```
License MIT.
