// TODO: switch to a better command line flag package
// TODO: implement more GNU ls options
// TODO: fix sorting
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

//%ls_colors = (
//	'README$'        => 11,
//	'Makefile$'      => $c[15],
//	'(=:.+)?\..*rc'  => $c[3],
//);

var lsColorSuffix = make(map[string]string)
var lsColorTypes = [...]string{
	"\033[",  // lc: Left of color sequence
	"m",      // rc: Right of color sequence
	"",       // ec: End color (replaces lc+no+rc)
	"0",      // rs: Reset to ordinary colors
	"",       // no: Normal
	"",       // fi: File: default
	"01;34",  // di: Directory: bright blue
	"01;36",  // ln: Symlink: bright cyan
	"33",     // pi: Pipe: yellow/brown
	"01;35",  // so: Socket: bright magenta
	"01;33",  // bd: Block device: bright yellow
	"01;33",  // cd: Char device: bright yellow
	"",       // mi: Missing file: undefined
	"",       // or: Orphaned symlink: undefined
	"01;32",  // ex: Executable: bright green
	"01;35",  // do: Door: bright magenta
	"37;41",  // su: setuid: white on red
	"30;43",  // sg: setgid: black on yellow
	"37;44",  // st: sticky: black on blue
	"34;42",  // ow: other-writable: blue on green
	"30;42",  // tw: ow w/ sticky: black on green
	"30;41",  // ca: black on red
	"",       // mh: disabled by default
	"\033[K", // cl: clear to end of line
}

func isColored(t indicator) bool {
	return lsColorTypes[t] != ""
}

func parseLSColor() {
	lsc := os.Getenv("LS_COLORS")
	var eq bool
	var kb, ke int
	for i := 0; i < len(lsc); i++ {
		b := lsc[i]
		if eq {
			if b == ':' {
				if lsc[kb] == '*' {
					lsColorSuffix[lsc[kb+1:ke]] = lsc[ke+1 : i]
				} else {
					k := lsc[kb:ke]
					fail := true
					for in, s := range indicatorNames {
						if s == k {
							lsColorTypes[in] = lsc[ke+1 : i]
							fail = false
						}
					}
					if fail {
						fmt.Printf("Unrecognized key: %s\n", k)
					}
				}
				i++
				kb = i
				eq = false
			}
		} else {
			if b == '=' {
				ke = i
				eq = true
			}
		}
	}
	//for k, v := range lsColorSuffix {
	//	fmt.Printf("%s = \033[%s;1m%s\033[0m\n", k, v, v)
	//}
}

type fileInfos []os.FileInfo

func (p fileInfos) Len() int { return len(p) }
func (p fileInfos) Less(i, j int) bool {
	mi, mj := p[i].IsDir(), p[j].IsDir()
	if mi == mj {
		return p[i].Name() < p[j].Name()
	}
	return mi
}
func (p fileInfos) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func main() {
	parseLSColor()
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	fis := make([]fileInfos, 0, len(args))
	for _, name := range args {
		nfis, err := ls(name)
		if err != nil {
			fmt.Println(err)
		} else {
			snfis := fileInfos(nfis)
			sort.Sort(snfis)
			fis = append(fis, snfis)
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
		name = ESC + cc + "m" + name + cEnd
	}
	if linkname != "" {
		lc := color(linkname, t)
		name = name + cSymDelim +
			ESC + "38;5;8;3m" +
			ESC + lc + "m" +
			linkname + cEnd
	}
	if *classify || *fclass {
		switch {
		case mode.IsDir():
			return name + "/"
		case t == TypeExec && !*fclass:
			return name + "*"
		case t == TypeFifo:
			return name + "|"
		case t == TypeSock:
			return name + "="
		}
	}
	return name
}

func color(name string, in indicator) string {
	if in == TypeFile {
		for i := 0; i < len(name); i++ {
			if name[i] == '.' {
				if v, ok := lsColorSuffix[name[i:]]; ok {
					return v
				}
			}
		}
	}
	return lsColorTypes[in]
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
