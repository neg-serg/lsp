// TODO: switch to a better command line flag package
// TODO: implement more GNU ls options
// TODO: fix sorting

//%ls_colors = (
//	'README$'        => 11,
//	'Makefile$'      => $c[15],
//	'(=:.+)?\..*rc'  => $c[3],
//);

package main

import (
	"fmt"
	"os"
	"sort"

	flag "github.com/neeee/pflag"
)

var (
	all      = flag.BoolP("all", "a", false, "show all")
	classify = flag.BoolP("classify", "F", false, "append indicator")
	fclass   = flag.Bool("file-type", false, "append indicators except *")
	_        = flag.BoolP("list", "l", false, "noop")
	_        = flag.BoolP("human-readable", "h", false, "noop")
)

func init() {
	flag.BoolVarP(all, "almost-all", "A", false, "show all")
}

const (
	cLeftCol  = "\033[38;5;0m" + "├"
	cRightCol = "\033[38;5;0m" + "┤" + cEnd
	cCol      = "\033[38;5;0m" + "│" + cEnd
	cSymDelim = " " + "\033[38;5;9m" + "→" + cEnd + " "
)

type fileList []os.FileInfo

func (p fileList) Len() int { return len(p) }
func (p fileList) Less(a, b int) bool {
	aF, bF := p[a], p[b]
	aD, bD := aF.IsDir(), bF.IsDir()
	if aD != bD {
		return aD
	}
	sA, sB := aF.Name(), bF.Name()
	return filevercmp(sA, sB) < 0
}
func (p fileList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

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
			snfis := fileList(nfis)
			sort.Sort(snfis)
			fis = append(fis, nfis)
		}
	}

	for _, fis := range fis {
		for _, f := range fis {
			fmt.Println(cLeftCol +
				strmode(f.Mode()) +
				cRightCol +
				reltime(f.ModTime()) +
				cCol +
				size(f.Size()) +
				cCol +
				name(f))
		}
	}
}

func name(f os.FileInfo) string {
	var l os.FileInfo
	linkok := true
	linkname := ""
	mode := f.Mode()
	if f.Mode()&os.ModeSymlink == os.ModeSymlink {
		var err error
		linkname, err = os.Readlink(f.Name())
		if err != nil {
			linkok = false
		} else {
			l, err = os.Stat(linkname)
			if err != nil {
				linkok = false
			} else {
				mode = l.Mode()
			}
		}
	}

	t := colorType(mode, linkok)
	cc := color(f.Name(), t)
	name := f.Name()
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
	if *classify || *fclass {
		switch {
		case mode.IsDir():
			return name + "/"
		case t == typeExec && !*fclass:
			return name + "*"
		case t == typeFifo:
			return name + "|"
		case t == typeSock:
			return name + "="
		}
	}
	return name
}

//ffis := make([]os.FileInfo, 0, len(fis))
//for _, fi := range fis {
//	if !strings.HasPrefix(fi.Name(), ".") {
//		ffis = append(ffis, fi)
//	}
//}

func ls(name string) ([]os.FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fis, err := f.Readdir(0)
		if *all {
			return fis, err
		}
		filtered := make([]os.FileInfo, 0, len(fis))
		for _, fi := range fis {
			if name := fi.Name(); len(name) > 0 && name[0] == '.' {
				continue
			}
			filtered = append(filtered, fi)
		}
		return filtered, err
	}
	return []os.FileInfo{fi}, nil
}
