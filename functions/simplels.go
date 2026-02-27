package functions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type color string
const (
	Blue color = "\x1b[34m"
	Green color = "\x1b[32m"
	Default color = "\x1b[0m"
)

func handleFile(w io.Writer, file string, color bool) error {
	if color {
		Default.colorPrint(w, file)
	} else {
		w.Write([]byte(file))
	}
	return nil
}

func printDir(w io.Writer, dir string, color bool) {
	if color {
		Blue.colorPrint(w, dir)
	} else {
		w.Write([]byte(dir))
	}
}

func printExec(w io.Writer, exec string, color bool) {
	if color {
		Green.colorPrint(w, exec)
	} else {
		w.Write([]byte(exec))
	}
}

/* from what I understand, ls (at least on the department machines)
sorts on lowercase comparisons. It also seems to ignore punctuation. */
func sort1(entries []os.DirEntry) ([]os.DirEntry, error) {
	slices.SortFunc(entries, func(a, b os.DirEntry) int {
			name1, name2 := a.Name(), b.Name()
			name1, name2 = strings.ReplaceAll(name1, ".", ""), strings.ReplaceAll(name2, ".", "")
			
			alow, blow := strings.ToLower(name1), strings.ToLower(name2)
			if alow != blow {
				if alow < blow { return -1 }
				return 1
			}
			
			return 0
		})

	return entries, nil
}

func handleDirectory(w io.Writer, dir string, color bool, header bool) error {

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	entries = dirFilter(entries)
	entries, err = sort1(entries)

	if header {
		w.Write([]byte(dir + ":\n"))
	}

	for i := 0; i < len(entries); i++ {
		cur := entries[i]
		curPath := filepath.Join(dir, cur.Name())

		info, err := os.Lstat(curPath)
		mode := info.Mode()
		if err != nil {
			return err
		}

		if info.IsDir() {
			printDir(w, cur.Name(), color)
		} else if mode.IsRegular() && (mode & 0111) != 0 {
			printExec(w, cur.Name(), color)
		} else {
			handleFile(w, cur.Name(), color)
		}

		if i != len(entries) - 1 {w.Write([]byte("\n"))}
	}

	return nil
}

func SimpleLS(w io.Writer, files []string, dirs []string, useColor bool) {
	header := len(dirs) + len(files) > 1

	for i := 0; i < len(files); i++ {

		info, err := os.Lstat(files[i])
		mode := info.Mode()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if mode.IsRegular() && (mode & 0111) != 0 {
			printExec(w, files[i], useColor)
			w.Write([]byte("\n"))
			continue
		}

		err = handleFile(w, files[i], useColor)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		w.Write([]byte("\n"))
	}

	if len(files) > 0 && len(dirs) > 0 {w.Write([]byte("\n"))}

	for i := 0; i < len(dirs); i++ {
		err := handleDirectory(w, dirs[i], useColor, header)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		w.Write([]byte("\n"))
		
		if i != len(dirs) - 1 {
			w.Write([]byte("\n"))
		}
	}
}