package functions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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

func indexMap(entries []os.DirEntry) map[int]int {

	var sorted []string
	for i := 0; i < len(entries); i++ {
		sorted = append(sorted, entries[i].Name())
	}

	original := sorted
	sort.Slice(sorted, func(i, j int) bool {
    	return strings.ToLower(sorted[i]) < strings.ToLower(sorted[j])
	})


	idx := make(map[int]int)
	for i := 0; i < len(original); i++ {
		idx[i] = slices.Index(sorted, original[i])
	}

	return idx
}

func transformEntries(entries []os.DirEntry, idx map[int]int) []os.DirEntry {
	out := make([]os.DirEntry, len(entries))
	for i := 0; i < len(entries); i++ {
		out[i] = entries[idx[i]]
	}

	return out
}

func sort1(entries []os.DirEntry) ([]os.DirEntry, error) {
	var out []os.DirEntry

	for i := 0; i < len(entries); i++ {
		cur := entries[i]
		out = append(out, cur)
	}

	out = transformEntries(out, indexMap(out))

	return out, nil
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