// Package sse writes Server-Sent Events in Datastar's wire format.
//
// Datastar (https://data-star.dev) defines exactly two event types:
//
//   - datastar-patch-elements: morph or replace HTML on the page
//   - datastar-patch-signals:  update client-side reactive signals
//
// This package is a hand-rolled minimal writer — small enough to read
// straight through, no SDK dependency.
package sse

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PatchMode mirrors Datastar's element-patch modes.
type PatchMode string

const (
	ModeOuter   PatchMode = "outer"   // morph outerHTML (default; finds target by id)
	ModeInner   PatchMode = "inner"   // morph innerHTML
	ModeReplace PatchMode = "replace" // hard replace, no morph
	ModePrepend PatchMode = "prepend"
	ModeAppend  PatchMode = "append"
	ModeBefore  PatchMode = "before"
	ModeAfter   PatchMode = "after"
	ModeRemove  PatchMode = "remove"
)

// Writer streams Datastar SSE events on a single HTTP response.
//
// Header setup is deferred until the first write so callers can still
// short-circuit with a non-SSE error response if needed.
type Writer struct {
	w        http.ResponseWriter
	flusher  http.Flusher
	prepared bool
}

// New wraps an http.ResponseWriter for SSE output.
func New(w http.ResponseWriter) *Writer {
	f, _ := w.(http.Flusher)
	return &Writer{w: w, flusher: f}
}

func (sw *Writer) prepare() {
	if sw.prepared {
		return
	}
	h := sw.w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no") // disable nginx buffering
	sw.prepared = true
}

// PatchElements emits one datastar-patch-elements event.
//
// selector is optional when mode is Outer/Replace and elements contains an
// element with an id matching the target. mode == "" defaults to Outer.
func (sw *Writer) PatchElements(selector string, mode PatchMode, elements string) error {
	sw.prepare()

	var b strings.Builder
	b.WriteString("event: datastar-patch-elements\n")
	if selector != "" {
		fmt.Fprintf(&b, "data: selector %s\n", selector)
	}
	if mode != "" {
		fmt.Fprintf(&b, "data: mode %s\n", mode)
	}
	for _, line := range strings.Split(elements, "\n") {
		fmt.Fprintf(&b, "data: elements %s\n", line)
	}
	b.WriteString("\n")

	if _, err := io.WriteString(sw.w, b.String()); err != nil {
		return err
	}
	if sw.flusher != nil {
		sw.flusher.Flush()
	}
	return nil
}
