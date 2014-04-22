package main

func size(w writer, size int64) {
	const (
		nnull  = "0123456789"
		nblank = " 123456789"
	)

	var (
		idx int
		dec int64
		buf [3]byte
		fmt string
	)

	for size > 999 {
		if size > 1024 {
			idx++
		}
		dec = size
		size /= 1024
	}
	// round <1.05 && >=0.95 to 1
	if dec < 1024/100 && dec >= 95*1024/100 {
		size = 1
		idx++
	}
	if dec < 105*1024/100 && dec >= 1024/100 {
		size = 1
	}
	dec = (dec*1000/1024 + 50) / 100 % 10

	fmt = nblank
	if size >= 10 || idx == 0 {
		// value is >= 10: use "123M', " 12M" formats
		buf[0] = fmt[size/100]
		if size/100 != 0 {
			fmt = nnull
		}
		buf[1] = fmt[size/10%10]
		buf[2] = byte('0' + size%10)
	} else {
		// value is < 10: use "9.2M" format
		buf[0] = byte('0' + size)
		buf[1] = '.'
		buf[2] = byte('0' + dec)
	}

	w.Write(cSize)
	w.Write(buf[:])
	w.Write(cSizes[idx])
}
