package logger

import (
	"io"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/httplog/v2"
)

func NewLogger(w io.Writer) *httplog.Logger {

	logger := httplog.NewLogger("server", httplog.Options{
		Writer:           w,
		Concise:          true,
		MessageFieldName: "message",
	})
	return logger
}

// so gopls does not remove it from my stuff CONSTANTLY
func SpewDump(a ...interface{}) {
	spew.Dump(a)
}
