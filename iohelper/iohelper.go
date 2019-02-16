// Package iohelper provides extra interfaces than package io.
package iohelper

import (
	"errors"
	"io"
)

// NewSectionWriter returns a SectionWriter that writes to w
// starting at offset off and stops with io.ErrShortWrite after n bytes.
func NewSectionWriter(w io.WriterAt, off int64, n int64) *SectionWriter {
	return &SectionWriter{w, off, off, off + n}
}

// SectionWriter implements Write, Seek, and WriteAt on a section
// of an underlying io.WriterAt.
type SectionWriter struct {
	w     io.WriterAt
	base  int64
	off   int64
	limit int64
}

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
func (s *SectionWriter) Size() int64 { return s.limit - s.base }
