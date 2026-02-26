package functions

import (
	"os"
	"regexp"
)

func dirFilter(entries []os.DirEntry) []os.DirEntry {
	var out []os.DirEntry

	for i := 0; i < len(entries); i++ {
		cur := entries[i]

		match, err := regexp.MatchString("^[\\.].*", cur.Name())
		if err != nil {
			os.Exit(2)
		}

		if !match {
			out = append(out, cur)
		}
	}

	return out
}