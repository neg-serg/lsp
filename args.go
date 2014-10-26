package main

import (
	"fmt"
	"os"
)

const usage = `Usage: lsp [option ...] [file ...]
  -C when  Use colours (never, always or auto)
  -F       Append file type indicator
  -a       Show all files
  -c       Use ctime
  -S       Sort by size
  -r       Reverse sort
  -t       Sort by time
  -h       Show this help`

type useColor uint8

const (
	colorNever useColor = iota
	colorAlways
	colorAuto
)

var opts = struct {
	all      bool
	classify bool
	color    useColor
	ctime    bool
	reverse  bool
	sorter   sortFunc
	rest     []string
}{
	color:  colorAuto,
	sorter: sortByVer,
	rest:   make([]string, 0, len(os.Args[1:])),
}

func parseArgs() {
	for a := os.Args[1:]; len(a) != 0; a = a[1:] {
		s := a[0]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			opts.rest = append(opts.rest, s)
			continue
		}
		if s[1] == '-' && len(s) == 2 { // "--" ends opts
			opts.rest = append(opts.rest, a[1:]...)
			break
		}
		for j := 1; j < len(s); j++ {
			f := s[j]
			switch f {
			case 'a':
				opts.all = true
			case 'F':
				opts.classify = true
			case 'c':
				opts.ctime = true
			case 'r':
				opts.reverse = true
			case 't':
				opts.sorter = sortByTime
			case 'S':
				opts.sorter = sortBySize
			case 'C':
				// XXX: ugly
				var use string
				if j != len(s)-1 {
					use = s[j+1:]
					j = len(s)
				} else {
					if len(a) == 1 {
						fmt.Fprintln(os.Stderr, "option '-C' needs an argument")
						fmt.Fprintln(os.Stderr, usage)
						os.Exit(1)
					}
					use = a[1]
					a = a[1:]
				}
				switch use {
				case "never":
					opts.color = colorNever
				case "always":
					opts.color = colorAlways
				case "auto":
					opts.color = colorAuto
				default:
					fmt.Fprintf(os.Stderr, "invalid argument to option '-C': %q\n", use)
					fmt.Fprintln(os.Stderr, usage)
					os.Exit(1)
				}
			case 'h':
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "unsupported option '%c'\n", f)
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(1)
			}
		}
	}
	if len(opts.rest) == 0 {
		opts.rest = []string{"."}
	}
}
