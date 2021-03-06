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

func errp(a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], fmt.Sprint(a))
}

func errf(f string, a ...interface{}) { fmt.Fprintf(os.Stderr, f, a...) }

func main() {
	parseArgs()
	colorize := (opts.color == colorAuto && isTTY(os.Stdout)) || opts.color == colorAlways
	if colorize {
		parseLSColor()
	}
	b := bufio.NewWriter(os.Stdout)
	for _, fname := range opts.rest {
		nfis, err := ls(fname)
		if err != nil {
			errp(err)
			continue
		}

		sorter := opts.sorter(nfis)
		if opts.reverse {
			sorter = sort.Reverse(sorter)
		}
		sort.Sort(sorter)
		if colorize {
			for _, f := range nfis {
				b.Write(nCol)
				strmode(b, f.mode)
				b.WriteByte(' ')
				b.Write(nCol)
				reltime(b, f.time)
				b.WriteByte(' ')
				b.Write(nCol)
				size(b, f.size)
				b.Write(cCol)
				name(b, f)
				b.WriteByte('\n')
			}
		} else {
			for _, f := range nfis {
				strmodeNoColor(b, f.mode)
				reltimeNoColor(b, f.time)
				b.Write(nCol)
				sizeNoColor(b, f.size)
				b.Write(nCol)
				nameNoColor(b, f)
				b.WriteByte('\n')
			}
		}
	}
	b.Flush()
}

func nameNoColor(b writer, f *fileInfo) {
	var t indicator
	if f.linkname.str != "" {
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

	b.WriteString(f.name.str)
	if f.linkname.str != "" {
		b.Write(nSymDelim)
		b.WriteString(f.linkname.str)
	}

	if opts.classify {
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

func name(b writer, f *fileInfo) {
	var t indicator
	if f.linkname.str != "" {
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
	c := color(f.name.str, t)
	if c == "" {
		b.WriteString(f.name.str)
	} else {
		b.Write(cESC)
		b.WriteString(c)
		b.WriteByte('m')
		b.WriteString(f.name.str)
		b.Write(cEnd)
	}

	if f.linkname.str != "" {
		var lnt indicator
		if !f.linkok {
			lnt = typeMissing
		} else {
			lnt = colorType(f.linkmode)
		}
		lc := color(f.linkname.str, lnt)
		if lc == "" {
			b.WriteString(f.linkname.str)
		} else {
			b.Write(cSymDelim)
			b.Write(cESC)
			b.WriteString("38;5;8;3m")
			b.Write(cESC)
			b.WriteString(lc)
			b.WriteByte('m')
			b.WriteString(f.linkname.str)
			b.Write(cEnd)
		}
	}

	if opts.classify {
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
