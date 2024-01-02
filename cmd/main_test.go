package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	go func() {
		var b bytes.Buffer
		cmd.Stdout = &b
		cmd.Stderr = &b
		err := cmd.Run()
		if err != nil {
			fmt.Printf("cmd.Run() failed: %s %s\n",
				err.Error(), b.String())
			panic(err)
		}
	}()
	var resp *http.Response
	for i := 0; i < 30; i++ {
		var err error
		resp, err = http.Get("http://localhost:1234/testdata/test.css.gz")
		if errors.Is(err, syscall.ECONNREFUSED) {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if err != nil {
			t.Fatalf("http.Get failed %s", err.Error())
		}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.Readall failed %s", err.Error())
	}
	body := string(b)
	assertEqual(t, "body {\n    margin: 25px;\n}", string(body))
	assertEqual(t, "text/css", resp.Header.Get("Content-Type"))
	assertEqual(t, "", resp.Header.Get("Content-Encoding"))
	assertEqual(t, true, resp.Uncompressed)
}

func assertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("expected (%+v) is not equal to actual (%+v)", expected, actual)
}
