package main

import (
	"syscall"
)

func gettime(st *syscall.Stat_t) int64 {
	if args.ctime {
		return int64(st.Ctim.Sec)*1e9 + int64(st.Ctim.Nsec)
	}
	return int64(st.Mtim.Sec)*1e9 + int64(st.Mtim.Nsec)
}

