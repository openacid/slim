// Package iohelper provides extra interfaces than package io.
package iohelper

import (
	"errors"
	"io"
)

const (
	maxOffset int64 = 0x7fffffffffffffff
)

// AtToWriter convert a WriterAt to a Writer
//
// Since 0.1.6
func AtToWriter(w io.WriterAt, offset int64) io.Writer {
	return NewSectionWriter(w, offset, maxOffset-offset)
}

// AtToReader convert a ReaderAt to a Reader
//
// Since 0.1.6
func AtToReader(r io.ReaderAt, offset int64) io.Reader {
	return io.NewSectionReader(r, offset, maxOffset-offset)
}

// NewSectionWriter returns a SectionWriter that writes to w
// starting at offset off and stops with io.ErrShortWrite after n bytes.
//
// Since 0.1.6
func NewSectionWriter(w io.WriterAt, off int64, n int64) *SectionWriter {
	return &SectionWriter{w, off, off, off + n}
}

// SectionWriter implements Write, Seek, and WriteAt on a section
// of an underlying io.WriterAt.
//
// Since 0.1.6
type SectionWriter struct {
	w     io.WriterAt
	base  int64
	off   int64
	limit int64
}

// Write buf "p".
//
// Since 0.1.6
func (s *SectionWriter) Write(p []byte) (n int, err error) {
	if s.off >= s.limit {
		return 0, io.ErrShortWrite
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
		err = io.ErrShortWrite
	}
	n, err2 := s.w.WriteAt(p, s.off)
	s.off += int64(n)

	if err2 != nil {
		err = err2
	}

	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

// Seek seeks to relative position by offset.
//
// Since 0.1.6
func (s *SectionWriter) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errWhence
	case io.SeekStart:
		offset += s.base
	case io.SeekCurrent:
		offset += s.off
	case io.SeekEnd:
		offset += s.limit
	}
	if offset < s.base {
		return 0, errOffset
	}
	s.off = offset
	return offset - s.base, nil
}

// WriteAt write buf p at relative position off.
//
// Since 0.1.6
func (s *SectionWriter) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off >= s.limit-s.base {
		return 0, io.ErrShortWrite
	}
	off += s.base
	if max := s.limit - off; int64(len(p)) > max {
		p = p[0:max]
		n, err = s.w.WriteAt(p, off)
		if err == nil {
			err = io.ErrShortWrite
		}
		return n, err
	}
	return s.w.WriteAt(p, off)
}

// Size returns the size of the section in bytes.
//
// Since 0.1.6
func (s *SectionWriter) Size() int64 { return s.limit - s.base }
