package main

const (
	sESC = "\033["     // Start escape sequence
	sEnd = sESC + "0m" // End escape sequence
)

var (
	// Column delimiters
	nCol = []byte(sESC + "38;5;235m" + "▏")

	// Symlink -> symlink target
	nSymDelim = []byte(" " + "→" + " ")
)

var (
	nNone = []byte("—") // Nothing else applies

	nRead  = []byte("r") // Readable
	nWrite = []byte("w") // Writeable
	nExec  = []byte("x") // Executable

	nDir   = []byte("d") // Directory
	nChar  = []byte("c") // Character device
	nBlock = []byte("b") // Block device
	nFifo  = []byte("p") // FIFO
	nLink  = []byte("l") // Symlink

	nSock    = []byte("s") // Socket
	nUID     = []byte("S") // SUID
	nUIDExec = []byte("s") // SUID and executable
	nSticky  = []byte("t") // Sticky
	nStickyO = []byte("T") // Sticky, writeable by others
)

// Relative times
var (
	nSecond = []byte{}
	nMinute = []byte{}
	nHour   = []byte{}
	nDay    = []byte{}
	nWeek   = []byte{}
	nMonth  = []byte{}
	nYear   = []byte{}
)

// Number part of size
var nSize = []byte{}

var nSizes = [...][]byte{
	[]byte("B"), // Byte
	[]byte("K"), // Kibibyte
	[]byte("M"), // Mebibyte
	[]byte("G"), // Gibibyte
	[]byte("T"), // Tebibyte
	[]byte("P"), // Pebibyte
	[]byte("E"), // Exbibyte
}

var (
	cESC = []byte(sESC)
	cEnd = []byte(sEnd)

	// Column delimiters
	cCol = []byte(sESC + "38;5;235m" + "│")

	// Symlink -> symlink target
	cSymDelim = []byte(" " + sESC + "38;5;8m" + "→" + sEnd + " ")
)

var (
	cNone = []byte(sESC + "38;5;240m" + "-") // Nothing else applies

	cRead  = []byte(sESC + "38;5;2m" + "r") // Readable
	cWrite = []byte(sESC + "38;5;7m" + "w") // Writeable
	cExec  = []byte(sESC + "38;5;1m" + "x") // Executable

	cDir   = []byte(sESC + "38;5;2m" + "d" + sEnd)     // Directory
	cChar  = []byte(sESC + "0m" + "c")                 // Character device
	cBlock = []byte(sESC + "0m" + "b")                 // Block device
	cFifo  = []byte(sESC + "0m" + "p")                 // FIFO
	cLink  = []byte(sESC + "38;5;220;1m" + "l" + sEnd) // Symlink

	cSock    = []byte(sESC + "38;5;161m" + "s")          // Socket
	cUID     = []byte(sESC + "38;5;220m" + "S")          // SUID
	cUIDExec = []byte(sESC + "38;5;161m" + "s")          // SUID and executable
	cSticky  = []byte(sESC + "38;5;220m" + "t")          // Sticky
	cStickyO = []byte(sESC + "38;5;220;1m" + "T" + sEnd) // Sticky, writeable by others
)

// Colours of relative times
var (
	cSecond = []byte(sESC + "38;5;4m")
	cMinute = []byte(sESC + "38;5;4m")
	cHour   = []byte(sESC + "00;38;5;12m")
	cDay    = []byte(sESC + "01;38;5;12m")
	cWeek   = []byte(sESC + "38;5;7m")
	cMonth  = []byte(sESC + "38;5;7m")
	cYear   = []byte(sESC + "38;5;8m")
)

// Number part of size
var cSize = []byte(sESC + "38;5;7m")

var cSizes = [...][]byte{
	[]byte(sESC + "38;5;7;1m" + "B" + sEnd),  // Byte
	[]byte(sESC + "38;5;2;1m" + "K" + sEnd),  // Kibibyte
	[]byte(sESC + "38;5;14;1m" + "M" + sEnd), // Mebibyte
	[]byte(sESC + "38;5;12;1m" + "G" + sEnd), // Gibibyte
	[]byte(sEnd + "T"),                       // Tebibyte
	[]byte(sEnd + "P"),                       // Pebibyte
	[]byte(sEnd + "E"),                       // Exbibyte
}
