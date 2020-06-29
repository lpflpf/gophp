go-serialize 
====

[![Build Status](https://travis-ci.com/lpflpf/gophp.svg?branch=master)](https://travis-ci.com/lpflpf/gophp)


Golang's implementation for PHP's function serialize() and unSerialize()


## Install / Update

```
go get -u github.com/lpflpf/gophp
```

## Example

```golang
package main

import (
	"fmt"

	"github.com/lpflpf/gophp"
)

func main() {

	str := `a:1:{s:3:"php";s:24:"世界上最好的语言";}`

	// unserialize() in php
	out, _ := gophp.UnMarshal([]byte(str))

	fmt.Println(out) //map[php:世界上最好的语言]

	// serialize() in php
	jsonbyte, _ := gophp.Marshal(out)

	fmt.Println(string(jsonbyte)) // a:1:{s:3:"php";s:24:"世界上最好的语言";}

}
```

