package contextio_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dikaeinstein/contextio"
)

func setupReader(ctx context.Context) io.Reader {
	r := strings.NewReader("hello")
	return contextio.NewReader(ctx, r)
}

func TestReaderCtx(t *testing.T) {
	rCtx := setupReader(context.Background())
	p := make([]byte, 5)

	n, err := rCtx.Read(p)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("5 bytes read expected; got %d", n)
	}
	if string(p) != "hello" {
		t.Error("Bad content")
	}

	p = make([]byte, 5)
	ctx, cancel := context.WithCancel(context.Background())
	rCtx = setupReader(ctx)
	n, err = rCtx.Read(p)

	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("5 bytes read expected; got %d", n)
	}
	if string(p) != "hello" {
		t.Error("Bad content")
	}
	cancel()
	n, err = rCtx.Read(p)
	if err != nil && err != context.Canceled {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("0 bytes read expected; got %d", n)
	}
}

func TestWriterCtx(t *testing.T) {
	var buf bytes.Buffer
	w := contextio.NewWriter(context.Background(), &buf)
	n, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("5 bytes written expected; got %d", n)
	}
	if buf.String() != "hello" {
		t.Error("Bad content")
	}

	// reset buffer
	buf.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	w = contextio.NewWriter(ctx, &buf)
	n, err = w.Write([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Errorf("5 bytes written expected; got %d", n)
	}
	if buf.String() != "hello" {
		t.Error("Bad content")
	}

	cancel()
	n, err = w.Write([]byte("hello"))
	if err != context.Canceled {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("0 bytes written expected; got %d", n)
	}
}

func TestCloserCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cCtx := contextio.NewCloser(ctx, ioutil.NopCloser(strings.NewReader("hello")))
	err := cCtx.Close()
	if err != nil {
		t.Fatal(err)
	}

	cancel()
	err = cCtx.Close()
	if err != nil && err != context.Canceled {
		t.Fatal(err)
	}
}
