package main

const (
	cESC = "\033["     // Start escape sequence
	cEnd = cESC + "0m" // End escape sequence
)

const (
	cCol      = " "                      // Column delimiters
	cSymDelim = " " + cESC + "38;5;9m" + // Symlink -> symlink target
		"→" + cEnd + " "
)

const (
	cNone = cESC + "38;5;0m" + "—" // Nothing else applies

	cRead  = cESC + "38;5;2m" + "r"   // Readable
	cWrite = cESC + "38;5;216m" + "w" // Writeable
	cExec  = cESC + "38;5;131m" + "x" // Executable

	cDir   = cESC + "38;5;2;1m" + "d" + cEnd   // Directory
	cChar  = cESC + "0m" + "c"                 // Character device
	cBlock = cESC + "0m" + "b"                 // Block device
	cFifo  = cESC + "0m" + "p"                 // FIFO
	cLink  = cESC + "38;5;220;1m" + "l" + cEnd // Symlink

	cSock    = cESC + "38;5;161m" + "s"          // Socket
	cUID     = cESC + "38;5;220m" + "S"          // SUID
	cUIDExec = cESC + "38;5;161m" + "s"          // SUID and executable
	cSticky  = cESC + "38;5;220m" + "t"          // Sticky
	cStickyO = cESC + "38;5;220;1m" + "T" + cEnd // Sticky, writeable by others
)

// Colours of relative times
const (
	cSecond = cESC + "38;5;12m"
	cMinute = cESC + "38;5;9m"
	cHour   = cESC + "38;5;1m"
	cDay    = cESC + "38;5;8m"
	cWeek   = cESC + "38;5;237m"
	cMonth  = cESC + "38;5;237m"
	cYear   = cESC + "38;5;0m"
)

// Number part of size
const cSize = cESC + "38;5;216m"

var cSizes = [...]string{
	cESC + "38;5;7;1m" + "B" + cEnd,  // Byte
	cESC + "38;5;2;1m" + "K" + cEnd,  // Kibibyte
	cESC + "38;5;14;1m" + "M" + cEnd, // Mebibyte
	cESC + "38;5;12;1m" + "G" + cEnd, // Gibibyte
	cEnd + "T",                       // Tebibyte
	cEnd + "P",                       // Pebibyte
	cEnd + "E",                       // Exbibyte
}
