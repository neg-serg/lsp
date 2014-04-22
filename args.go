package main

import (
	"fmt"
	"os"
)

const usage = `Usage: lsp -[aAFcrtS] [file ...]
  -a  Show all files
  -c  Use ctime
  -F  Append file type indicator
  -r  Reverse sort
  -t  Sort by time
  -S  Sort by size
  -h  Show this help`

type lsargs struct {
	all      bool
	classify bool
	ctime    bool
	reverse  bool
	sorter   sortFunc
	rest     []string
}

type unsupportedError struct {
	Flag rune
}

func (e *unsupportedError) Error() string {
	return fmt.Sprintf("unsupported argument '%c'", e.Flag)
}

func parseArgs(args []string) *lsargs {
	la := lsargs{
		all:      false,
		classify: false,
		ctime:    false,
		reverse:  false,
		sorter:   sortByVer,
		rest:     make([]string, 0, len(args)),
	}
	for _, s := range args {
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			la.rest = append(la.rest, s)
			continue
		}
		if s[1] == '-' && len(s) == 2 { // "--" ends args
			break
		}
		for i := 1; i < len(s); i++ {
			f := s[i]
			switch f {
			case 'a':
				la.all = true
			case 'F':
				la.classify = true
			case 'c':
				la.ctime = true
			case 'r':
				la.reverse = true
			case 't':
				la.sorter = sortByTime
			case 'S':
				la.sorter = sortBySize
			case 'h':
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "unsupported argument '%c'\n", f)
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(1)
			}
		}
	}
	if len(la.rest) == 0 {
		la.rest = []string{"."}
	}
	return &la
}
