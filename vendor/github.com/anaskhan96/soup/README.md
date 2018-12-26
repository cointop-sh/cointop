# soup
[![Build Status](https://travis-ci.org/anaskhan96/soup.svg?branch=master)](https://travis-ci.org/anaskhan96/soup)
[![GoDoc](https://godoc.org/github.com/anaskhan96/soup?status.svg)](https://godoc.org/github.com/anaskhan96/soup)
[![Go Report Card](https://goreportcard.com/badge/github.com/anaskhan96/soup)](https://goreportcard.com/report/github.com/anaskhan96/soup)

**Web Scraper in Go, similar to BeautifulSoup**

*soup* is a small web scraper package for Go, with its interface highly similar to that of BeautifulSoup.

Exported variables and functions implemented till now :
```go
var Headers map[string]string // Set headers as a map of key-value pairs, an alternative to calling Header() individually
var Cookies map[string]string // Set cookies as a map of key-value  pairs, an alternative to calling Cookie() individually
func Get(string) (string,error){} // Takes the url as an argument, returns HTML string
func GetWithClient(string, *http.Client){} // Takes the url and a custom HTTP client as arguments, returns HTML string
func Header(string, string){} // Takes key,value pair to set as headers for the HTTP request made in Get()
func Cookie(string, string){} // Takes key, value pair to set as cookies to be sent with the HTTP request in Get()
func HTMLParse(string) Root {} // Takes the HTML string as an argument, returns a pointer to the DOM constructed
func Find([]string) Root {} // Element tag,(attribute key-value pair) as argument, pointer to first occurence returned
func FindAll([]string) []Root {} // Same as Find(), but pointers to all occurrences returned
func FindStrict([]string) Root {} //  Element tag,(attribute key-value pair) as argument, pointer to first occurence returned with exact matching values
func FindAllStrict([]string) []Root {} // Same as FindStrict(), but pointers to all occurrences returned
func FindNextSibling() Root {} // Pointer to the next sibling of the Element in the DOM returned
func FindNextElementSibling() Root {} // Pointer to the next element sibling of the Element in the DOM returned
func FindPrevSibling() Root {} // Pointer to the previous sibling of the Element in the DOM returned
func FindPrevElementSibling() Root {} // Pointer to the previous element sibling of the Element in the DOM returned
func Children() []Root {} // Find all direct children of this DOM element
func Attrs() map[string]string {} // Map returned with all the attributes of the Element as lookup to their respective values
func Text() string {} // Full text inside a non-nested tag returned, first half returned in a non-nested one
func FullText() string {} // Full text inside a nested/non-nested tag returned
func SetDebug(bool) {} // Sets the debug mode to true or false; false by default
```

`Root` is a struct, containing three fields :
* `Pointer` containing the pointer to the current html node
* `NodeValue` containing the current html node's value, i.e. the tag name for an ElementNode, or the text in case of a TextNode
* `Error` containing an error if one occurrs, else `nil` is returned.

## Installation
Install the package using the command
```bash
go get github.com/anaskhan96/soup
```

## Example
An example code is given below to scrape the "Comics I Enjoy" part (text and its links) from [xkcd](https://xkcd.com).

[More Examples](https://github.com/anaskhan96/soup/tree/master/examples)
```go
package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"os"
)

func main() {
	resp, err := soup.Get("https://xkcd.com")
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	links := doc.Find("div", "id", "comicLinks").FindAll("a")
	for _, link := range links {
		fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
	}
}
```

## Contributions
This package was developed in my free time. However, contributions from everybody in the community are welcome, to make it a better web scraper. If you think there should be a particular feature or function included in the package, feel free to open up a new issue or pull request.
