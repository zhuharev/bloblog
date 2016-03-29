package bloblog

import (
	"fmt"
	"io"
	"os"
)

var (
	DefaultIndexSize int64 = 8 * 100000 //MB
)

type BlobLog struct {
	f *os.File

	indexSize int64
}

func New(fpath string) (*BlobLog, error) {
	f, e := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0666)
	if e != nil {
		return nil, e
	}
	stat, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if stat.Size() == 0 {
		e = f.Truncate(DefaultIndexSize)
		if e != nil {
			return nil, e
		}
	}
	bl := new(BlobLog)
	bl.f = f
	bl.indexSize = DefaultIndexSize
	return bl, nil
}

func (bl *BlobLog) LastInserId() (int64, error) {
	var res = make([]byte, 8)
	_, e := bl.f.ReadAt(res, 0)
	if e != nil {
		if e != io.EOF {
			return 0, e
		}
	}
	return b2i(res), nil
}

func (bl *BlobLog) Prepare(size int64) (int64, error) {
	lid, e := bl.LastInserId()
	if e != nil {
		return 0, e
	}
	newId := lid + 1

	stat, e := bl.f.Stat()
	if e != nil {
		return 0, e
	}
	e = bl.f.Truncate(stat.Size() + size)
	if e != nil {
		return 0, e
	}
	off := 8 + (lid * 16)
	_, e = bl.f.Seek(off, 0)
	if e != nil {
		return 0, e
	}
	_, e = bl.f.Write(i2b(stat.Size()))
	if e != nil {
		return 0, e
	}
	_, e = bl.f.Write(i2b(size))
	if e != nil {
		return 0, e
	}

	_, e = bl.f.WriteAt(i2b(newId), 0)
	if e != nil {
		return 0, e
	}
	return newId, nil
}

func (bl *BlobLog) GetMeta(id int64) (int64, int64, error) {
	off := (id * 16) - 8
	var res = make([]byte, 16)
	n, e := bl.f.ReadAt(res, off)
	if e != nil {
		return 0, 0, e
	}
	if n != 16 {
		return 0, 0, fmt.Errorf("n!=16")
	}
	return b2i(res[:8]), b2i(res[8:]), nil
}

func (bl *BlobLog) Write(id int64, in []byte) error {
	offset, size, e := bl.GetMeta(id)
	if e != nil {
		return e
	}
	if int64(len(in)) != size {
		return fmt.Errorf("Size input not equal prepared size got %d, need %d", len(in), size)
	}
	_, e = bl.f.Seek(offset, 0)
	if e != nil {
		return e
	}
	_, e = bl.f.Write(in)
	return e
}

func (bl *BlobLog) Insert(data []byte) (id int64, e error) {
	id, e = bl.Prepare(len(data))
	if e != nil {
		return 0, e
	}
	e = bl.Write(id, data)
	return
}

func (bl *BlobLog) Get(id int64) ([]byte, error) {
	offset, size, e := bl.GetMeta(id)
	if e != nil {
		return nil, e
	}
	res := make([]byte, size)
	_, e = bl.f.ReadAt(res, offset)
	if e != nil {
		return nil, e
	}
	return res, nil
}

type limitWriter struct {
	limit int64
	cur   int64
	io.Writer
}

func (lw *limitWriter) Write(in []byte) (int, error) {
	return 0, nil
}
