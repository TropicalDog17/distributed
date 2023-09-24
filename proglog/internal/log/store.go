package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append the data p to the file buffer wrapped by store s.
// Length of the data is also included at the beginning of the appended p.
// Return number of bytes written, position where the store hold the data in the file and error if any
// Record simple form: <DATA_SIZE: 8 bytes><DATA: DATA_SIZE bytes>
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Last index
	pos = s.size
	// We write the length of the record so we know how many bytes to read
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	// len(p) = 8 bytes
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

// Read record starting at pos.
// Return the byte data of the record and error if any
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	// Read the size of data
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	// Proceed to read the data after knowing the data size
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// Read len(p) bytes into p beginning, starting at the off offset in the record file.
// Acts as a wrapper of File.ReadAt()
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

//
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}
