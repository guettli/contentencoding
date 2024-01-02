package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/guettli/contentencoding"
)

func usage() {
	fmt.Printf(`%s [address (default: localhost:1234)] [directory (default: .)]
Serve files from the local directory via http.
Gzipped files like foo.css.gz will be served with Content-Encoding gzip and the 
appropriate Content-Type.

Example, serve the directory "static" on localhost with port 8080:
 %s localhost:8080 static

`, os.Args[0], os.Args[0])
}

func main() {
	if len(os.Args) > 1 && os.Args[1][0] == '-' {
		usage()
		os.Exit(0)
	}
	addr := "localhost:1234"
	directory := "."
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	if len(os.Args) > 2 {
		directory = os.Args[2]
		_, err := os.Stat(directory)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	if len(os.Args) > 3 {
		usage()
		os.Exit(0)
	}
	http.Handle("/", contentencoding.FileServer(http.Dir(directory)))
	fmt.Printf("Listening on http://%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
	os.Exit(1)
}
