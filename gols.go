package main

import (
	"fmt"
	"io"
	"ls-part-one-AbElo-77/functions"
	"os"
	"regexp"
	"sort"
)

func filterFile(file string) bool {
	_, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return false
		}

	return true
}

func parseArgs(args []string) ([]string, []string, error) {
	var files []string
	var dirs []string

	for i := 0; i < len(args); i++ {
		match, err := regexp.MatchString(".+[\\.].+", args[i])
		if err != nil {
			return nil, nil, err
		}

		if match {
			ok := filterFile(args[i])
			
			if ok {
				files = append(files, args[i])
			}
		} else {
			info, err := os.Lstat(args[i])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			mode := info.Mode()

			if mode.IsRegular() && (mode & 0111) != 0 {
				files = append(files, args[i])
			}

			ok := mode.IsDir()
			if ok {
				dirs = append(dirs, args[i])
			}

		}
	}

	return files, dirs, nil
	
}

func handleArgs(files []string, dirs []string) {

	var writer io.Writer = os.Stdout

	isTerminal := functions.IsTerminal(os.Stdout)
	functions.SimpleLS(writer, files, dirs, isTerminal)
}

func main() {

	args := os.Args
	if len(args) == 1 {
		curdir, err := os.Getwd()
		if err != nil {
			return 
		}

		var dirlist []string
		dirlist = append(dirlist, curdir)

		handleArgs(nil, dirlist)
		return
	}

	args = args[1:]

	files, dirs, err := parseArgs(args)
	if err != nil {
		os.Exit(2)
	}

	sort.Strings(files)
	sort.Strings(dirs)

	handleArgs(files, dirs)
}