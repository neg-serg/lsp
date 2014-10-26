package main

import (
	"syscall"
)

func gettime(st *syscall.Stat_t) int64 {
	if opts.ctime {
		return int64(st.Ctimespec.Sec)*1e9 + int64(st.Ctimespec.Nsec)
	}
	return int64(st.Mtimespec.Sec)*1e9 + int64(st.Mtimespec.Nsec)
}
