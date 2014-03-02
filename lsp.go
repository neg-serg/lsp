// TODO: switch to a better command line flag package
// TODO: implement more GNU ls options

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
	"syscall"
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

//	all      = flag.BoolP("all", "a", false, "show all")
//	classify = flag.BoolP("classify", "F", false, "append indicator")
//	ctime    = flag.BoolP("ctime", "c", false, "ctime instead of modtime")
//	reverse  = flag.BoolP("reverse", "r", false, "reverse order while sorting")
//	sorttime = flag.BoolP("timesort", "t", false, "sort by time")
//	sortsize = flag.BoolP("sizesort", "S", false, "sort by size")
//	_        = flag.BoolP("list", "l", false, "noop")
//	_        = flag.BoolP("human-readable", "h", false, "noop")

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
		nfis, isdir, err := ls(fname)
		if err != nil {
			errp(err)
			continue
		}

		var cwd string
		if isdir {
			cwd, _ = os.Getwd()
			os.Chdir(fname)
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
			name(b, f)
			b.WriteByte('\n')
		}

		if isdir {
			os.Chdir(cwd)
		}
	}
	b.Flush()
}

func name(b writer, f *fileInfo) {
	var l *fileInfo
	linkok := true
	linkname := ""
	mode := f.mode
	if f.mode&syscall.S_IFMT == syscall.S_IFLNK {
		var err error
		linkname, err = readlink(f.name)
		if err != nil {
			linkok = false
		} else {
			l, err = stat(linkname)
			if err != nil {
				linkok = false
			} else {
				mode = l.mode
			}
		}
	}

	t := colorType(mode, linkok)
	cc := color(f.name, t)
	if cc == "" {
		b.WriteString(f.name)
	} else {
		b.WriteString(cESC)
		b.WriteString(cc)
		b.WriteByte('m')
		b.WriteString(f.name)
		b.WriteString(cEnd)
	}
	if linkname != "" {
		b.WriteString(cSymDelim + cESC + "38;5;8;3m" + cESC)
		b.WriteString(color(linkname, t))
		b.WriteByte('m')
		b.WriteString(linkname)
		b.WriteString(cEnd)
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
