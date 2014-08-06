package main

const (
	sESC = "\033["     // Start escape sequence
	sEnd = sESC + "0m" // End escape sequence
)

var (
	cESC = []byte(sESC)
	cEnd = []byte(sEnd)

	// Column delimiters
	cCol = []byte(" ")

	// Symlink -> symlink target
	cSymDelim = []byte(" " + sESC + "38;5;8m" + "→" + sEnd + " ")
)

var (
	cNone = []byte(sESC + "38;5;0m" + "—") // Nothing else applies

	cRead  = []byte(sESC + "38;5;2m" + "r")   // Readable
	cWrite = []byte(sESC + "38;5;216m" + "w") // Writeable
	cExec  = []byte(sESC + "38;5;131m" + "x") // Executable

	cDir   = []byte(sESC + "38;5;2;1m" + "d" + sEnd)   // Directory
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
	cSecond = []byte(sESC + "38;5;12m")
	cMinute = []byte(sESC + "38;5;9m")
	cHour   = []byte(sESC + "38;5;1m")
	cDay    = []byte(sESC + "38;5;8m")
	cWeek   = []byte(sESC + "38;5;237m")
	cMonth  = []byte(sESC + "38;5;237m")
	cYear   = []byte(sESC + "38;5;0m")
)

// Number part of size
var cSize = []byte(sESC + "38;5;216m")

var cSizes = [...][]byte{
	[]byte(sESC + "38;5;7;1m" + "B" + sEnd),  // Byte
	[]byte(sESC + "38;5;2;1m" + "K" + sEnd),  // Kibibyte
	[]byte(sESC + "38;5;14;1m" + "M" + sEnd), // Mebibyte
	[]byte(sESC + "38;5;12;1m" + "G" + sEnd), // Gibibyte
	[]byte(sEnd + "T"),                       // Tebibyte
	[]byte(sEnd + "P"),                       // Pebibyte
	[]byte(sEnd + "E"),                       // Exbibyte
}
