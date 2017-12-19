# GopherNet

[![Build Status](https://travis-ci.org/pallavagarwal07/gophernet.svg?branch=master)](https://travis-ci.org/pallavagarwal07/gophernet)

Go library for network calls that uses net/http or jquery conditionally
compiled with either the Go compiler or gopherjs.

## Utility

To test this out, create a new folder with a main.go file:

```
package main

import "github.com/pallavagarwal07/gophernet"

func main() {
	URL := `https://raw.githubusercontent.com/pallavagarwal07/gophernet/master/README.md`
	out, err := gophernet.Get(URL, nil)
	if err != nil {
		println(err)
		panic("Error")
	}
	println(string(out))
}
```

Now, run the file using

```
go run main.go
```

The output should contain the contents of README.md of this repository.
Now let's try this with GopherJS. Install GopherJS:

```
go get -u github.com/gopherjs/gopherjs
```

Now, in the same directory as main.go, run `gopherjs serve`.  If everything is
working correctly, this should run a server at `localhost:8080`. Open it in
browser, and see the developer's console. It should have printed out the
contents of the README as before.

If you are not familiar with GopherJS, what happens here is that the program
you wrote has been compiled to javascript and the server serves a blank HTML
page with just that javascript file in the `<head>` section. GopherNet provides
a way to do HTTP requests from `go` that get compiled to javascript (using XHR)
on using gopherjs by using proper build tags inside the source file(s).
