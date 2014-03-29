// TODO: switch to a better command line flag package
// TODO: (maybe) implement more GNU ls options
// TODO: more flexible colours

//%ls_colors = (
//	'README$'        => 11,
//	'Makefile$'      => $c[15],
//	'(=:.+)?\..*rc'  => $c[3],
//);

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const shortUsage = "Usage: %s -[aAFcrtS] [file ...]\n"
const longUsage = `
  -a, -A  Show all files
  -F      Append file type indicator
  -c      Use ctime
  -r      Reverse sort
  -t      Sort by time
  -S      Sort by size
`

type lsargs struct {
	all      bool
	classify bool
	ctime    bool
	reverse  bool
	sorttime bool
	sortsize bool
	rest     []string
}

var args = &lsargs{}

func parseArgs() {
	for _, s := range os.Args[1:] {
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			args.rest = append(args.rest, s)
			continue
		}
		if s[1] == '-' && len(s) == 2 { // "--" ends args
			break
		}
		for _, f := range s[1:] {
			switch f {
			case 'A', 'a':
				args.all = true
			case 'F':
				args.classify = true
			case 'c':
				args.ctime = true
			case 'r':
				args.reverse = true
			case 't':
				args.sorttime = true
				args.sortsize = false
			case 'S':
				args.sorttime = false
				args.sortsize = true
			case 'h':
				fmt.Printf(shortUsage, os.Args[0])
				fmt.Print(longUsage)
				os.Exit(0)
			default:
				errf("unsupported argument '%c'", f)
				os.Exit(1)
			}
		}
	}
	if len(args.rest) == 0 {
		args.rest = []string{"."}
	}
}

func errp(a ...interface{})           { fmt.Fprintln(os.Stderr, a...) }
func errf(f string, a ...interface{}) { fmt.Fprintf(os.Stderr, f, a...) }

func main() {
	var b = bufio.NewWriter(os.Stdout)
	parseArgs()
	parseLSColor()

	var sortFunc sortFunc
	if args.sorttime {
		sortFunc = sortByTime
	} else if args.sortsize {
		sortFunc = sortBySize
	} else {
		sortFunc = sortByVer
	}

	for _, fname := range args.rest {
		nfis, err := ls(fname)
		if err != nil {
			errp(err)
			continue
		}

		sorter := sortFunc(nfis)
		if args.reverse {
			sorter = sort.Reverse(sorter)
		}
		sort.Sort(sorter)
		for _, f := range nfis {
			strmode(b, f.mode)
			reltime(b, f.time)
			b.WriteString(cCol)
			size(b, f.size)
			b.WriteString(cCol)
			name(b, &f)
			b.WriteByte('\n')
		}
	}
	b.Flush()
}

func name(b writer, f *fileInfo) {
	var t indicator
	if f.linkname != "" {
		if !f.linkok {
			t = typeOrphan
		} else if colorSymTarget {
			t = colorType(f.linkmode)
		} else {
			t = colorType(f.mode)
		}
	} else {
		t = colorType(f.mode)
	}
	c := color(f.name, t)
	if c == "" {
		b.WriteString(f.name)
	} else {
		b.WriteString(cESC)
		b.WriteString(c)
		b.WriteByte('m')
		b.WriteString(f.name)
		b.WriteString(cEnd)
	}

	if f.linkname != "" {
		var lnt indicator
		if !f.linkok {
			lnt = typeMissing
		} else {
			lnt = colorType(f.linkmode)
		}
		lc := color(f.linkname, lnt)
		if lc == "" {
			b.WriteString(f.linkname)
		} else {
			b.WriteString(cSymDelim + cESC + "38;5;8;3m" + cESC)
			b.WriteString(color(f.linkname, lnt))
			b.WriteByte('m')
			b.WriteString(f.linkname)
			b.WriteString(cEnd)
		}
	}

	if args.classify {
		switch t {
		case typeDir:
			b.WriteByte('/')
		case typeExec:
			b.WriteByte('*')
		case typeFifo:
			b.WriteByte('|')
		case typeSock:
			b.WriteByte('=')
		}
	}
}
