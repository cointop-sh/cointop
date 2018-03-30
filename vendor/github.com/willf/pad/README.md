pad
-------------

[![Join the chat at https://gitter.im/willf/pad](https://badges.gitter.im/willf/pad.svg)](https://gitter.im/willf/pad?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

[![Build Status](https://travis-ci.org/willf/pad.svg?branch=master)](https://travis-ci.org/willf/pad)

A golang implementation of the [left-pad javascript library](https://www.npmjs.com/package/left-pad)

I was inspired by [Stew](https://twitter.com/StewOConnor)'s [`left-cats`](https://github.com/stew/left-cats), who was inspired by [this article](http://arstechnica.com/information-technology/2016/03/rage-quit-coder-unpublished-17-lines-of-javascript-and-broke-the-internet/), to port this to Go. 

This implementation will let you pad byte-strings and UTF-8 encoded strings

example usage:

```go
package main

import (
	"fmt"

	"github.com/willf/pad"
	padUtf8 "github.com/willf/pad/utf8"
)

func main() {
	fmt.Println(pad.Right("Hello", 20, "!"))
	fmt.Println(padUtf8.Left("Exit now", 20, "→"))
}
```

```bash
> go run example.go
Hello!!!!!!!!!!!!!!!
→→→→→→→→→→→→Exit now
```
