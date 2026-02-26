package functions

import (
	"io"
)

func (c color) colorPrint(w io.Writer, s string) {
	w.Write([]byte(c))
	w.Write([]byte(s))

	w.Write([]byte(Default))
}