// TODO: switch to a better command line flag package
// TODO: implement more GNU ls options

//%ls_colors = (
//	'README$'        => 11,
//	'Makefile$'      => $c[15],
//	'(=:.+)?\..*rc'  => $c[3],
//);

package main

import (
	"fmt"
	"sort"
	"syscall"

	flag "github.com/neeee/pflag"
)

var (
	all      = flag.BoolP("all", "a", false, "show all")
	classify = flag.BoolP("classify", "F", false, "append indicator")
	ctime    = flag.BoolP("ctime", "c", false, "ctime instead of modtime")
	_        = flag.BoolP("list", "l", false, "noop")
	_        = flag.BoolP("human-readable", "h", false, "noop")
)

func init() {
	flag.BoolVarP(all, "almost-all", "A", false, "show all")
}

func main() {
	parseLSColor()
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	fis := make([]fileList, 0, len(args))
	for _, name := range args {
		nfis, err := ls(name)
		if err != nil {
			fmt.Println(err)
		} else {
			sort.Sort(byVer(nfis))
			fis = append(fis, nfis)
		}
	}

	for _, fis := range fis {
		for _, f := range fis {
			fmt.Println(cLeftCol +
				strmode(f.mode) +
				cRightCol +
				reltime(f.time) +
				cCol +
				size(f.size) +
				cCol +
				name(f))
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
