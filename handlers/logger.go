package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

var (
	middlewareDefaultLogger func(next http.Handler) http.Handler

	// Normal colors
	nBlack   = []byte{'\033', '[', '3', '0', 'm'}
	nRed     = []byte{'\033', '[', '3', '1', 'm'}
	nGreen   = []byte{'\033', '[', '3', '2', 'm'}
	nYellow  = []byte{'\033', '[', '3', '3', 'm'}
	nBlue    = []byte{'\033', '[', '3', '4', 'm'}
	nMagenta = []byte{'\033', '[', '3', '5', 'm'}
	nCyan    = []byte{'\033', '[', '3', '6', 'm'}
	nWhite   = []byte{'\033', '[', '3', '7', 'm'}

	// Bright colors
	bBlack   = []byte{'\033', '[', '3', '0', ';', '1', 'm'}
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bYellow  = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

// DefaultLogFormatter is a simple logger that implements a LogFormatter.
type DefaultLogFormatter struct {
	UseColor bool
}

// NewLogEntry creates a new LogEntry for the request.
func (l *DefaultLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &defaultLogEntry{
		DefaultLogFormatter: l,
		request:             r,
		buf:                 &bytes.Buffer{},
		useColor:            l.UseColor,
	}

	cW(entry.buf, l.UseColor, bMagenta, "%s ", r.Method)
	cW(entry.buf, l.UseColor, nCyan, "%s %s", r.RequestURI, r.Proto)

	entry.buf.WriteString(" from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteRune(' ')
	return entry
}

type defaultLogEntry struct {
	*DefaultLogFormatter
	request  *http.Request
	buf      *bytes.Buffer
	useColor bool
}

func (l *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	switch {
	case status < 200:
		cW(l.buf, l.useColor, bBlue, "%03d", status)
	case status < 300:
		cW(l.buf, l.useColor, bGreen, "%03d", status)
	case status < 400:
		cW(l.buf, l.useColor, bCyan, "%03d", status)
	case status < 500:
		cW(l.buf, l.useColor, bYellow, "%03d", status)
	default:
		cW(l.buf, l.useColor, bRed, "%03d", status)
	}

	cW(l.buf, l.useColor, bBlue, " %dB", bytes)

	l.buf.WriteString(" in ")
	if elapsed < 50*time.Millisecond {
		cW(l.buf, l.useColor, nGreen, "%s", elapsed)
	} else if elapsed < 1*time.Second {
		cW(l.buf, l.useColor, nYellow, "%s", elapsed)
	} else {
		cW(l.buf, l.useColor, nRed, "%s", elapsed)
	}

	fmt.Fprintf(l.buf, " %+v\n", l.request.Header)
	log.Print(l.buf.String())
}

func (l *defaultLogEntry) Panic(v interface{}, stack []byte) {
	middleware.PrintPrettyStack(v)
}

// colorWrite
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if useColor {
		w.Write(reset)
	}
}

func init() {
	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	middlewareDefaultLogger = middleware.RequestLogger(&DefaultLogFormatter{UseColor: color})
}
