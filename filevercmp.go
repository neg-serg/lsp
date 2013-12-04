package main

// Helper functions

func isalpha(c byte) bool {
	return ((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z'))
}

func isalnum(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z')
}

func isdigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

// Match a file suffix defined by this regular expression:
// /(\.[A-Za-z~][A-Za-z0-9~]*)*$/
// Scan the string *STR and return a pointer to the matching suffix, or
// NULL if not found.  Upon return, *STR points to terminating NUL.
func suffixIndex(s string) int {
	read_alpha := false
	match := -1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if read_alpha {
			read_alpha = false
			if !isalpha(c) && c != '~' {
				match = -1
			}
		} else if '.' == c {
			read_alpha = true
			if match == -1 {
				match = i
			}
		} else if !isalnum(c) && c != '~' {
			match = -1
		}
	}
	return match
}

// verrevcmp helper function
func order(c byte) int {
	if isdigit(c) {
		return 0
	} else if isalpha(c) {
		return int(c)
	} else if c == '~' {
		return -1
	} else {
		return int(c) + 256
	}
}

// Slightly modified verrevcmp function from dpkg
func verrevcmp(a, b string) int {
	ai, bi := 0, 0
	for ai < len(a) || bi < len(b) {
		first_diff := 0
		for (ai < len(a) && !isdigit(a[ai])) || (bi < len(b) && !isdigit(b[bi])) {
			var ac, bc int
			if ai < len(a) {
				ac = order(a[ai])
			}
			if bi < len(b) {
				bc = order(b[bi])
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
		for ai < len(a) && isdigit(a[ai]) && bi < len(b) && isdigit(b[bi]) {
			if first_diff == 0 {
				first_diff = int(a[ai]) - int(b[bi])
			}
			ai++
			bi++
		}

		if ai < len(a) && isdigit(a[ai]) {
			return 1
		}
		if bi < len(b) && isdigit(b[bi]) {
			return -1
		}
		if first_diff != 0 {
			return first_diff
		}
	}

	return 0
}

// Compare version strings s1 and s2
func filevercmp(s1, s2 string) int {
	var result int

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
	switch {
	case (s1[0] == '.' && s2[0] != '.'):
		return -1
	case (s1[0] != '.' && s2[0] == '.'):
		return 1
	}
	if s1[0] == '.' && s2[0] == '.' {
		s1 = s1[1:]
		s2 = s2[1:]
	}

	// file suffixes
	s1_suffix_index := suffixIndex(s1)
	s2_suffix_index := suffixIndex(s2)

	var s1_cut, s2_cut = s1, s2
	var s1_suffix, s2_suffix string
	if s1_suffix_index != -1 {
		s1_cut = s1[:s1_suffix_index]
		s1_suffix = s1[s1_suffix_index:]
	}
	if s2_suffix_index != -1 {
		s2_cut = s2[:s2_suffix_index]
		s2_suffix = s2[s2_suffix_index:]
	}

	// restore file suffixes if strings are identical after "cut"
	if s1_cut == s2_cut {
		result = verrevcmp(s1_suffix, s2_suffix)
	} else {
		result = verrevcmp(s1_cut, s2_cut)
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
