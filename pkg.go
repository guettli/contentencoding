package contentencoding

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// ContentEncodingFileServer is a thin wrapper around
// http.FileServer(). It servers zipped files with the appropriate
// Content-Encoding http response header, so that you can use zipped
// files directly.
// Example: You want to use that in html:
//
//	<script src="/static/htmx.min.js.gz"></script>
//
// Go code:
// http.Handle("/static/", http.StripPrefix("/static/",
//
//	contentencodingfileserver.FileServer(http.Dir("./static")))))
func FileServer(root http.FileSystem) http.Handler {
	return SetContentEncodingHandler(http.FileServer(root))
}

type contentEncodingReponseWriter struct {
	wrapped http.ResponseWriter
	request *http.Request
}

func (w *contentEncodingReponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

func (w *contentEncodingReponseWriter) Write(b []byte) (int, error) {
	return w.wrapped.Write(b)
}

func (w *contentEncodingReponseWriter) WriteHeader(statusCode int) {
	w.rewriteHeader(statusCode)
	w.wrapped.WriteHeader(statusCode)
}

func (w *contentEncodingReponseWriter) rewriteHeader(statusCode int) {
	r := w.request
	if statusCode != 200 {
		return
	}
	if !strings.HasSuffix(r.URL.Path, ".gz") {
		return
	}
	if w.Header().Get("Content-Type") != "application/gzip" {
		return
	}
	ext := filepath.Ext(strings.TrimSuffix(r.URL.Path, ".gz"))
	if ext == "" {
		return
	}
	ct := w.detectContentType(ext)
	if ct == "" {
		return
	}
	// Up to now the code was read-only.
	// Now we are going to re-write the headers.
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", ct)
}

func (w *contentEncodingReponseWriter) detectContentType(ext string) string {
	switch ext {
	case ".css":
		return "text/css"
	case ".js":
		return "text/javascript"
	case ".html":
		return "text/html"
	default:
		return mime.TypeByExtension(ext)
	}
}

// SetContentEncodingHandler wraps a http.Handler and
// sets Content-Encoding to "gzip" and Content-Type to the
// appropriate type if request.URL.Path ends with .css.gz, .js.gz, or similar.
func SetContentEncodingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&contentEncodingReponseWriter{wrapped: w, request: r}, r)
	})
}
