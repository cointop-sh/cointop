// Copyright 2017 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mp3

import (
	"io"
)

type source struct {
	reader io.ReadCloser
	buf    []byte
	pos    int64
}

func (s *source) Seek(position int64, whence int) (int64, error) {
	seeker, ok := s.reader.(io.Seeker)
	if !ok {
		panic("mp3: source must be io.Seeker")
	}
	s.buf = nil
	n, err := seeker.Seek(position, whence)
	if err != nil {
		return 0, err
	}
	s.pos = n
	return n, nil
}

func (s *source) Close() error {
	s.buf = nil
	return s.reader.Close()
}

func (s *source) skipTags() error {
	buf := make([]byte, 3)
	if _, err := s.ReadFull(buf); err != nil {
		return err
	}
	switch string(buf) {
	case "TAG":
		buf := make([]byte, 125)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}

	case "ID3":
		// Skip version (2 bytes) and flag (1 byte)
		buf := make([]byte, 3)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}

		buf = make([]byte, 4)
		n, err := s.ReadFull(buf)
		if err != nil {
			return err
		}
		if n != 4 {
			return nil
		}
		size := (uint32(buf[0]) << 21) | (uint32(buf[1]) << 14) |
			(uint32(buf[2]) << 7) | uint32(buf[3])
		buf = make([]byte, size)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}

	default:
		s.Unread(buf)
	}

	return nil
}

func (s *source) rewind() error {
	if _, err := s.Seek(0, io.SeekStart); err != nil {
		return err
	}
	s.pos = 0
	s.buf = nil
	return nil
}

func (s *source) Unread(buf []byte) {
	s.buf = append(s.buf, buf...)
	s.pos -= int64(len(buf))
}

func (s *source) ReadFull(buf []byte) (int, error) {
	read := 0
	if s.buf != nil {
		read = copy(buf, s.buf)
		if len(s.buf) > read {
			s.buf = s.buf[read:]
		} else {
			s.buf = nil
		}
		if len(buf) == read {
			return read, nil
		}
	}

	n, err := io.ReadFull(s.reader, buf[read:])
	if err != nil {
		// Allow if all data can't be read. This is common.
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
	}
	s.pos += int64(n)
	return n + read, err
}
