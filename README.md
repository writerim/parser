# parser


## Example
```golang
package main

import (
  "github.com/parser"
)

func main() {

  p := parser.New()
  err := p.ConnectDb("127.0.0.1", 3306, "root", "", "parser")
  if err != nil {
    return
  }

  go p.StartWebServer(8080)
  p.StartDeamon()

}


```