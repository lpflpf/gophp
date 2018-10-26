# gophp
Golang implementation for PHP's functions

## Install / Update

```
go get -u github.com/techleeone/gophp
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/techleeone/gophp/serialize"
)

func main() {

	str := `a:1:{s:3:"php";s:24:"世界上最好的语言";}`

	// unserialize() in php
	out, _ := serialize.UnMarshal([]byte(str))

	fmt.Println(out) //map[php:世界上最好的语言]

	// serialize() in php
	jsonbyte, _ := serialize.Marshal(out)

	fmt.Println(string(jsonbyte)) // a:1:{s:3:"php";s:24:"世界上最好的语言";}

}
```
