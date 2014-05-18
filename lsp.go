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
	parseLSColor()
	var b = bufio.NewWriter(os.Stdout)
	for _, fname := range args.rest {
		nfis, err := ls(fname)
		if err != nil {
			errp(err)
			continue
		}

		sorter := args.sorter(nfis)
		if args.reverse {
			sorter = sort.Reverse(sorter)
		}
		sort.Sort(sorter)
		for _, f := range nfis {
			strmode(b, f.mode)
			reltime(b, f.time)
			b.Write(cCol)
			size(b, f.size)
			b.Write(cCol)
			name(b, f)
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
		b.Write(cESC)
		b.WriteString(c)
		b.WriteByte('m')
		b.WriteString(f.name)
		b.Write(cEnd)
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
			b.Write(cSymDelim)
			b.Write(cESC)
			b.WriteString("38;5;8;3m")
			b.Write(cESC)
			b.WriteString(color(f.linkname, lnt))
			b.WriteByte('m')
			b.WriteString(f.linkname)
			b.Write(cEnd)
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
