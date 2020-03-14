// Package contextio wraps io based operations to make them context-aware.
package contextio

import (
	"context"
	"io"
)

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

// NewReader returns a context aware io.Reader that reads from r.
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	if r, ok := r.(*readerCtx); ok && ctx == r.ctx {
		return r
	}
	return readerCtx{
		ctx: ctx,
		r:   r,
	}
}

func (rc readerCtx) Read(p []byte) (n int, err error) {
	select {
	case <-rc.ctx.Done():
		return 0, rc.ctx.Err()
	default:
		return rc.r.Read(p)
	}
}

type writerCtx struct {
	ctx context.Context
	w   io.Writer
}

// NewWriter returns a context aware io.Writer that writes to w.
func NewWriter(ctx context.Context, w io.Writer) io.Writer {
	if w, ok := w.(*writerCtx); ok && ctx == w.ctx {
		return w
	}
	return writerCtx{
		ctx: ctx,
		w:   w,
	}
}

func (wc writerCtx) Write(p []byte) (n int, err error) {
	select {
	case <-wc.ctx.Done():
		return 0, wc.ctx.Err()
	default:
		return wc.w.Write(p)
	}
}

type closerCtx struct {
	ctx context.Context
	c   io.Closer
}

// NewCloser returns a context aware io.Closer that closes the reader or writer.
func NewCloser(ctx context.Context, c io.Closer) io.Closer {
	return closerCtx{
		ctx: ctx,
		c:   c,
	}
}

func (cc closerCtx) Close() error {
	select {
	case <-cc.ctx.Done():
		return cc.ctx.Err()
	default:
		return cc.c.Close()
	}
}
