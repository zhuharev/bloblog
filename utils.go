package bloblog

import (
	"encoding/binary"
)

func i2b(i int64) []byte {
	var res = make([]byte, 8)
	binary.BigEndian.PutUint64(res, uint64(i))
	return res
}

func b2i(b []byte) int64 {
	if b == nil || len(b) != 8 {
		return 0
	}
	ui := binary.BigEndian.Uint64(b)
	return int64(ui)
}
