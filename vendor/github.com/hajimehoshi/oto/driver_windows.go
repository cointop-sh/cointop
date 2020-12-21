// Copyright 2015 Hajime Hoshi
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

// +build !js

package oto

import (
	"errors"
	"runtime"
	"unsafe"
)

type header struct {
	buffer  []byte
	waveHdr *wavehdr
}

func newHeader(waveOut uintptr, bufferSize int) (*header, error) {
	h := &header{
		buffer: make([]byte, bufferSize),
	}
	h.waveHdr = &wavehdr{
		lpData:         uintptr(unsafe.Pointer(&h.buffer[0])),
		dwBufferLength: uint32(bufferSize),
	}
	if err := waveOutPrepareHeader(waveOut, h.waveHdr); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *header) Write(waveOut uintptr, data []byte) error {
	if len(data) != len(h.buffer) {
		return errors.New("oto: len(data) must equal to len(h.buffer)")
	}
	copy(h.buffer, data)
	if err := waveOutWrite(waveOut, h.waveHdr); err != nil {
		return err
	}
	return nil
}

type driver struct {
	out        uintptr
	headers    []*header
	tmp        []byte
	bufferSize int
}

func newDriver(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes int) (*driver, error) {
	numBlockAlign := channelNum * bitDepthInBytes
	f := &waveformatex{
		wFormatTag:      waveFormatPCM,
		nChannels:       uint16(channelNum),
		nSamplesPerSec:  uint32(sampleRate),
		nAvgBytesPerSec: uint32(sampleRate * numBlockAlign),
		wBitsPerSample:  uint16(bitDepthInBytes * 8),
		nBlockAlign:     uint16(numBlockAlign),
	}
	w, err := waveOutOpen(f)
	if err != nil {
		return nil, err
	}

	const numBufs = 2
	p := &driver{
		out:        w,
		headers:    make([]*header, numBufs),
		bufferSize: bufferSizeInBytes,
	}
	runtime.SetFinalizer(p, (*driver).Close)
	for i := range p.headers {
		var err error
		p.headers[i], err = newHeader(w, p.bufferSize)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *driver) TryWrite(data []byte) (int, error) {
	n := min(len(data), max(0, p.bufferSize-len(p.tmp)))
	p.tmp = append(p.tmp, data[:n]...)
	if len(p.tmp) < p.bufferSize {
		return n, nil
	}

	var headerToWrite *header
	for _, h := range p.headers {
		// TODO: Need to check WHDR_DONE?
		if h.waveHdr.dwFlags&whdrInqueue == 0 {
			headerToWrite = h
			break
		}
	}
	if headerToWrite == nil {
		return n, nil
	}

	if err := headerToWrite.Write(p.out, p.tmp); err != nil {
		// This error can happen when e.g. a new HDMI connection is detected (#51).
		const errorNotFound = 1168
		werr := err.(*winmmError)
		if werr.fname == "waveOutWrite" && werr.errno == errorNotFound {
			return 0, nil
		}
		return 0, err
	}

	p.tmp = nil
	return n, nil
}

func (p *driver) Close() error {
	runtime.SetFinalizer(p, nil)
	// TODO: Call waveOutUnprepareHeader here
	if err := waveOutClose(p.out); err != nil {
		return err
	}
	return nil
}
