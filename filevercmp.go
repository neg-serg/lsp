package main

type sufIndexed struct {
	str string
	idx int
}

func newSufIndexed(s string) *sufIndexed {
	if len(s) != 0 && s[0] == '.' {
		return &sufIndexed{s, suffixIndex(s[1:])}
	}
	return &sufIndexed{s, suffixIndex(s)}
}

func isAlpha(c byte) bool {
	return (((c) | 32) - 'a') < 26
}

func isAlnum(c byte) bool {
	return c-'0' < 10 || (((c)|32)-'a') < 26
}

func isDigit(c byte) bool {
	return c-'0' < 10
}

// Match a file suffix defined by this regular expression:
// /(\.[A-Za-z~][A-Za-z0-9~]*)*$/
// Scan the string *STR and return a pointer to the matching suffix, or
// NULL if not found.  Upon return, *STR points to terminating NUL.
func suffixIndex(s string) int {
	readAlphat := false
	matched := 0
	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		// Manual inlining helps a lot here
		if (((c)|32)-'a') < 26 || c == '~' { // isAlpha
			readAlphat = true
		} else if readAlphat && c == '.' {
			matched = j + 1
		} else if c-'0' < 10 { // isDigit
			readAlphat = false
		} else {
			break
		}
		j++
	}
	return len(s) - matched
}

// verrevcmp helper function
func sortOrder(c byte) int {
	if isAlpha(c) {
		return int(c)
	} else if isDigit(c) {
		return 0
	} else if c == '~' {
		return -1
	}
	return int(c) + 256
}

// Slightly modified verrevcmp function from dpkg
func verrevcmp(a, b string) int {
	ai, bi := 0, 0
	for ai < len(a) || bi < len(b) {
		firstDiff := 0
		for (ai < len(a) && !isDigit(a[ai])) ||
			(bi < len(b) && !isDigit(b[bi])) {
			var ac, bc int
			if ai < len(a) {
				ac = sortOrder(a[ai])
			}
			if bi < len(b) {
				bc = sortOrder(b[bi])
			}

			if ac != bc {
				return ac - bc
			}

			ai++
			bi++
		}
		for ai < len(a) && a[ai] == '0' {
			ai++
		}
		for bi < len(b) && b[bi] == '0' {
			bi++
		}
		for ai < len(a) && isDigit(a[ai]) &&
			bi < len(b) && isDigit(b[bi]) {
			if firstDiff == 0 {
				firstDiff = int(a[ai]) - int(b[bi])
			}
			ai++
			bi++
		}

		if ai < len(a) && isDigit(a[ai]) {
			return 1
		}
		if bi < len(b) && isDigit(b[bi]) {
			return -1
		}
		if firstDiff != 0 {
			return firstDiff
		}
	}

	return 0
}

// Compare version strings s1 and s2
func filevercmp(sf1, sf2 *sufIndexed) int {
	s1, s2 := sf1.str, sf2.str
	// easy comparison to see if strings are identical
	if s1 == s2 {
		return 0
	}

	// special handle for "", "." and ".."
	switch {
	case s1 == "":
		return -1
	case s2 == "":
		return 1
	case s1 == ".":
		return -1
	case s2 == ".":
		return 1
	case s1 == "..":
		return -1
	case s2 == "..":
		return 1
	}

	// special handle for other hidden files
	if s1[0] == '.' {
		if s2[0] == '.' {
			s1 = s1[1:]
			s2 = s2[1:]
		} else {
			return -1
		}
	} else if s2[0] == '.' {
		return 1
	}

	// file suffixes
	s1i := sf1.idx
	s2i := sf2.idx
	s1Cut, s1Suf := s1[:s1i], s1[s1i:]
	s2Cut, s2Suf := s2[:s2i], s2[s2i:]

	var result int
	// restore file suffixes if strings are identical after "cut"
	if s1Cut == s2Cut {
		result = verrevcmp(s1Suf, s2Suf)
	} else {
		result = verrevcmp(s1Cut, s2Cut)
	}

	if result == 0 {
		if s1 < s2 {
			return -1
		} else if s1 > s2 {
			return 1
		}
	}
	return result
}
