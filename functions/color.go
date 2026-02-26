package functions

import (
	"fmt"
	"io"
)

func (c color) colorPrint(w io.Writer, s string) {
	fmt.Fprint(w, c)
	fmt.Fprint(w, s)

	fmt.Fprint(w, Default)
}