# Serve gzipped files with Content-Encoding


Serve files from the local directory via http.

Gzipped files like foo.css.gz will be served with Content-Encoding gzip and the
appropriate Content-Type.

The Content-Type gets determined by the file extension with the help of [mime.TypeByExtension](https://pkg.go.dev/mime#TypeByExtension)

## Use case

If you want to use a file like `foo.js.gz`, then most web servers will serve the file with the Content-Type "application/gzip".

But this means, you can't use the file like in this html snippet:

```
<script src="/static/foo.js.gz"></script>
```

The error message of the browser:

> Uncaught SyntaxError: Invalid or unexpected token (at foo.js.gz:1:1)

The `contentencoding.FileServer` of this Go package serves the above file with these headers, so that the browser is able to consume these files.

> Content-Type: text/css
> 
> Content-Encoding: gzip

This way, the above html snippet just works, and the static gzipped files can be served without modification.


## Command-Line Usage:

```
❯ go run github.com/guettli/contentencoding/cmd@latest

Listening on http://localhost:1234
```

Arguments:

> [address (default: localhost:1234)] [directory (default: .)]


## Golang Package

If you want to serve the directory "static" under the URL prefix "/static", you can use this Go code:

```
import (
	"net/http"

	"github.com/guettli/contentencoding"
)

    ...
    http.Handle("/static/", http.StripPrefix("/static/",
        contentencoding.FileServer(http.Dir("./static"))))
```

## Feedback is welcome

Feel free to create an issue at Github to provide feedback.