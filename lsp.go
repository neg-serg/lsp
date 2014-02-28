// TODO: switch to a better command line flag package
// TODO: implement more GNU ls options

//%ls_colors = (
//	'README$'        => 11,
//	'Makefile$'      => $c[15],
//	'(=:.+)?\..*rc'  => $c[3],
//);

package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"syscall"

	flag "github.com/neeee/pflag"
)

var (
	all      = flag.BoolP("all", "a", false, "show all")
	classify = flag.BoolP("classify", "F", false, "append indicator")
	ctime    = flag.BoolP("ctime", "c", false, "ctime instead of modtime")
	reverse  = flag.BoolP("reverse", "r", false, "reverse order while sorting")
	sorttime = flag.BoolP("timesort", "t", false, "sort by time")
	sortsize = flag.BoolP("sizesort", "S", false, "sort by size")
	_        = flag.BoolP("list", "l", false, "noop")
	_        = flag.BoolP("human-readable", "h", false, "noop")
)

func init() {
	flag.BoolVarP(all, "almost-all", "A", false, "show all")
}

func main() {
	var b bytes.Buffer
	parseLSColor()
	flag.Parse()
	var sortFunc sortFunc
	if *sorttime {
		sortFunc = sortByTime
	} else if *sortsize {
		sortFunc = sortBySize
	} else {
		sortFunc = sortByVer
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, fname := range args {
		nfis, isdir, err := ls(fname)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var cwd string
		if isdir {
			cwd, _ = os.Getwd()
			os.Chdir(fname)
		}

		sorter := sortFunc(nfis)
		if *reverse {
			sorter = sort.Reverse(sorter)
		}
		sort.Sort(sorter)
		for _, f := range nfis {
			b.Write(strmode(f.mode))
			b.WriteString(reltime(f.time))
			b.Write([]byte(cCol))
			b.WriteString(size(f.size))
			b.Write([]byte(cCol))
			b.WriteString(name(f))
			b.WriteByte('\n')
		}
		b.WriteTo(os.Stdout)
		b.Reset()

		if isdir {
			os.Chdir(cwd)
		}
	}
}

func name(f *fileInfo) string {
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
	name := f.name
	if cc != "" {
		name = cESC + cc + "m" + name + cEnd
	}
	if linkname != "" {
		lc := color(linkname, t)
		name = name + cSymDelim +
			cESC + "38;5;8;3m" +
			cESC + lc + "m" +
			linkname + cEnd
	}
	if *classify {
		switch {
		case mode.isDir():
			return name + "/"
		case t == typeExec:
			return name + "*"
		case t == typeFifo:
			return name + "|"
		case t == typeSock:
			return name + "="
		}
	}
	return name
}
