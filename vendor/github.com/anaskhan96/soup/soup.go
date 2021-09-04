/* soup package implements a simple web scraper for Go,
keeping it as similar as possible to BeautifulSoup
*/

package soup

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Root is a structure containing a pointer to an html node, the node value, and an error variable to return an error if occurred
type Root struct {
	Pointer   *html.Node
	NodeValue string
	Error     error
}

var debug = false

// Headers contains all HTTP headers to send
var Headers = make(map[string]string)

// SetDebug sets the debug status
// Setting this to true causes the panics to be thrown and logged onto the console.
// Setting this to false causes the errors to be saved in the Error field in the returned struct.
func SetDebug(d bool) {
	debug = d
}

// Header sets a new HTTP header
func Header(n string, v string) {
	Headers[n] = v
}

// Get returns the HTML returned by the url in string
func Get(url string) (string, error) {
	defer catchPanic("Get()")
	// Init a new HTTP client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", errors.New("Couldn't perform GET request to " + url)
	}
	// Set headers
	for hName, hValue := range Headers {
		req.Header.Set(hName, hValue)
	}
	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + url)
		}
		return "", errors.New("Couldn't perform GET request to " + url)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			panic("Unable to read the response body")
		}
		return "", errors.New("Unable to read the response body")
	}
	return string(bytes), nil
}

// HTMLParse parses the HTML returning a start pointer to the DOM
func HTMLParse(s string) Root {
	defer catchPanic("HTMLParse()")
	r, err := html.Parse(strings.NewReader(s))
	if err != nil {
		if debug {
			panic("Unable to parse the HTML")
		}
		return Root{nil, "", errors.New("Unable to parse the HTML")}
	}
	for r.Type != html.ElementNode {
		switch r.Type {
		case html.DocumentNode:
			r = r.FirstChild
		case html.DoctypeNode:
			r = r.NextSibling
		case html.CommentNode:
			r = r.NextSibling
		}
	}
	return Root{r, r.Data, nil}
}

// Find finds the first occurrence of the given tag name,
// with or without attribute key and value specified,
// and returns a struct with a pointer to it
func (r Root) Find(args ...string) Root {
	defer catchPanic("Find()")
	temp, ok := findOnce(r.Pointer, args, false)
	if ok == false {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return Root{nil, "", errors.New("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")}
	}
	return Root{temp, temp.Data, nil}
}

// FindAll finds all occurrences of the given tag name,
// with or without key and value specified,
// and returns an array of structs, each having
// the respective pointers
func (r Root) FindAll(args ...string) []Root {
	defer catchPanic("FindAll()")
	temp := findAllofem(r.Pointer, args)
	if len(temp) == 0 {
		if debug {
			panic("Element `" + args[0] + "` with attributes `" + strings.Join(args[1:], " ") + "` not found")
		}
		return []Root{}
	}
	pointers := make([]Root, 0, 10)
	for i := 0; i < len(temp); i++ {
		pointers = append(pointers, Root{temp[i], temp[i].Data, nil})
	}
	return pointers
}

// FindNextSibling finds the next sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextSibling() Root {
	defer catchPanic("FindNextSibling()")
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		if debug {
			panic("No next sibling found")
		}
		return Root{nil, "", errors.New("No next sibling found")}
	}
	return Root{nextSibling, nextSibling.Data, nil}
}

// FindPrevSibling finds the previous sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevSibling() Root {
	defer catchPanic("FindPrevSibling()")
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		if debug {
			panic("No previous sibling found")
		}
		return Root{nil, "", errors.New("No previous sibling found")}
	}
	return Root{prevSibling, prevSibling.Data, nil}
}

// FindNextElementSibling finds the next element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindNextElementSibling() Root {
	defer catchPanic("FindNextElementSibling()")
	nextSibling := r.Pointer.NextSibling
	if nextSibling == nil {
		if debug {
			panic("No next element sibling found")
		}
		return Root{nil, "", errors.New("No next element sibling found")}
	}
	if nextSibling.Type == html.ElementNode {
		return Root{nextSibling, nextSibling.Data, nil}
	}
	p := Root{nextSibling, nextSibling.Data, nil}
	return p.FindNextElementSibling()
}

// FindPrevElementSibling finds the previous element sibling of the pointer in the DOM
// returning a struct with a pointer to it
func (r Root) FindPrevElementSibling() Root {
	defer catchPanic("FindPrevElementSibling()")
	prevSibling := r.Pointer.PrevSibling
	if prevSibling == nil {
		if debug {
			panic("No previous element sibling found")
		}
		return Root{nil, "", errors.New("No previous element sibling found")}
	}
	if prevSibling.Type == html.ElementNode {
		return Root{prevSibling, prevSibling.Data, nil}
	}
	p := Root{prevSibling, prevSibling.Data, nil}
	return p.FindPrevElementSibling()
}

// Attrs returns a map containing all attributes
func (r Root) Attrs() map[string]string {
	defer catchPanic("Attrs()")
	if r.Pointer.Type != html.ElementNode {
		if debug {
			panic("Not an ElementNode")
		}
		return nil
	}
	if len(r.Pointer.Attr) == 0 {
		return nil
	}
	return getKeyValue(r.Pointer.Attr)
}

// Text returns the string inside a non-nested element
func (r Root) Text() string {
	defer catchPanic("Text()")
	k := r.Pointer.FirstChild
checkNode:
	if k.Type != html.TextNode {
		k = k.NextSibling
		if k == nil {
			if debug {
				panic("No text node found")
			}
			return ""
		}
		goto checkNode
	}
	if k != nil {
		r, _ := regexp.Compile(`^\s+$`)
		if ok := r.MatchString(k.Data); ok {
			k = k.NextSibling
			if k == nil {
				if debug {
					panic("No text node found")
				}
				return ""
			}
			goto checkNode
		}
		return k.Data
	}
	return ""
}

// Using depth first search to find the first occurrence and return
func findOnce(n *html.Node, args []string, uni bool) (*html.Node, bool) {
	if uni == true {
		if n.Type == html.ElementNode && n.Data == args[0] {
			if len(args) > 1 && len(args) < 4 {
				for i := 0; i < len(n.Attr); i++ {
					if n.Attr[i].Key == args[1] && n.Attr[i].Val == args[2] {
						return n, true
					}
				}
			} else if len(args) == 1 {
				return n, true
			}
		}
	}
	uni = true
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p, q := findOnce(c, args, true)
		if q != false {
			return p, q
		}
	}
	return nil, false
}

// Using depth first search to find all occurrences and return
func findAllofem(n *html.Node, args []string) []*html.Node {
	var nodeLinks = make([]*html.Node, 0, 10)
	var f func(*html.Node, []string, bool)
	f = func(n *html.Node, args []string, uni bool) {
		if uni == true {
			if n.Data == args[0] {
				if len(args) > 1 && len(args) < 4 {
					for i := 0; i < len(n.Attr); i++ {
						if n.Attr[i].Key == args[1] && n.Attr[i].Val == args[2] {
							nodeLinks = append(nodeLinks, n)
						}
					}
				} else if len(args) == 1 {
					nodeLinks = append(nodeLinks, n)
				}
			}
		}
		uni = true
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, args, true)
		}
	}
	f(n, args, false)
	return nodeLinks
}

// Returns a key pair value (like a dictionary) for each attribute
func getKeyValue(attributes []html.Attribute) map[string]string {
	var keyvalues = make(map[string]string)
	for i := 0; i < len(attributes); i++ {
		_, exists := keyvalues[attributes[i].Key]
		if exists == false {
			keyvalues[attributes[i].Key] = attributes[i].Val
		}
	}
	return keyvalues
}

// Catch panics when they occur
func catchPanic(fnName string) {
	if r := recover(); r != nil {
		log.Println("Error occurred in", fnName, ":", r)
	}
}
