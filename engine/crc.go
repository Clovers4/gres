package engine

import (
	"bufio"
	"github.com/clovers4/gres/util"
	"hash"
	"hash/crc32"
)

type CRCWriter struct {
	w   *bufio.Writer
	crc hash.Hash32
}

func NewCRCWriter(w *bufio.Writer) *CRCWriter {
	return &CRCWriter{
		w:   w,
		crc: crc32.NewIEEE(),
	}
}

func (w *CRCWriter) Write(p []byte) (nn int, err error) {
	nn, err = w.crc.Write(p)
	if nn > 0 {
		if nn, err := w.w.Write(p[:nn]); err != nil {
			return nn, err
		}
	}
	return
}

func (w *CRCWriter) WriteCRC() error {
	crc := w.crc.Sum32()
	return util.Write(w.w, crc)
}

func (w *CRCWriter) Flush() error {
	return w.w.Flush()
}

type CRCReader struct {
	r   *bufio.Reader
	crc hash.Hash32
}

func NewCRCReader(r *bufio.Reader) *CRCReader {
	return &CRCReader{
		r:   r,
		crc: crc32.NewIEEE(),
	}
}

func (r *CRCReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		if n, err := r.crc.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (r *CRCReader) IsCRCEqual() (bool, error) {
	expect := r.crc.Sum32()

	var real uint32
	if err := util.Read(r, &real); err != nil {
		return false, err
	}
	return expect == real, nil
}
