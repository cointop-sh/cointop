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

// +build !js

package oto

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	winmm = windows.NewLazySystemDLL("winmm")
)

var (
	procWaveOutOpen          = winmm.NewProc("waveOutOpen")
	procWaveOutClose         = winmm.NewProc("waveOutClose")
	procWaveOutPrepareHeader = winmm.NewProc("waveOutPrepareHeader")
	procWaveOutWrite         = winmm.NewProc("waveOutWrite")
)

type wavehdr struct {
	lpData          uintptr
	dwBufferLength  uint32
	dwBytesRecorded uint32
	dwUser          uintptr
	dwFlags         uint32
	dwLoops         uint32
	lpNext          uintptr
	reserved        uintptr
}

type waveformatex struct {
	wFormatTag      uint16
	nChannels       uint16
	nSamplesPerSec  uint32
	nAvgBytesPerSec uint32
	nBlockAlign     uint16
	wBitsPerSample  uint16
	cbSize          uint16
}

const (
	waveFormatPCM = 1
	whdrInqueue   = 16
)

type mmresult uint

const (
	mmsyserrNoerror       mmresult = 0
	mmsyserrError         mmresult = 1
	mmsyserrBaddeviceid   mmresult = 2
	mmsyserrAllocated     mmresult = 4
	mmsyserrInvalidhandle mmresult = 5
	mmsyserrNodriver      mmresult = 6
	mmsyserrNomem         mmresult = 7
	waveerrBadformat      mmresult = 32
	waveerrStillplaying   mmresult = 33
	waveerrUnprepared     mmresult = 34
	waveerrSync           mmresult = 35
)

func (m mmresult) String() string {
	switch m {
	case mmsyserrNoerror:
		return "MMSYSERR_NOERROR"
	case mmsyserrError:
		return "MMSYSERR_ERROR"
	case mmsyserrBaddeviceid:
		return "MMSYSERR_BADDEVICEID"
	case mmsyserrAllocated:
		return "MMSYSERR_ALLOCATED"
	case mmsyserrInvalidhandle:
		return "MMSYSERR_INVALIDHANDLE"
	case mmsyserrNodriver:
		return "MMSYSERR_NODRIVER"
	case mmsyserrNomem:
		return "MMSYSERR_NOMEM"
	case waveerrBadformat:
		return "WAVEERR_BADFORMAT"
	case waveerrStillplaying:
		return "WAVEERR_STILLPLAYING"
	case waveerrUnprepared:
		return "WAVEERR_UNPREPARED"
	case waveerrSync:
		return "WAVEERR_SYNC"
	}
	return fmt.Sprintf("MMRESULT (%d)", m)
}

type winmmError struct {
	fname    string
	errno    windows.Errno
	mmresult mmresult
}

func (e *winmmError) Error() string {
	if e.errno != 0 {
		return fmt.Sprintf("winmm error at %s: Errno: %d", e.fname, e.errno)
	}
	if e.mmresult != mmsyserrNoerror {
		return fmt.Sprintf("winmm error at %s: %s", e.fname, e.mmresult)
	}
	return fmt.Sprintf("winmm error at %s", e.fname)
}

func waveOutOpen(f *waveformatex) (uintptr, error) {
	const (
		waveMapper   = 0xffffffff
		callbackNull = 0
	)
	var w uintptr
	r, _, e := procWaveOutOpen.Call(uintptr(unsafe.Pointer(&w)), waveMapper, uintptr(unsafe.Pointer(f)),
		0, 0, callbackNull)
	if e.(windows.Errno) != 0 {
		return 0, &winmmError{
			fname: "waveOutOpen",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return 0, &winmmError{
			fname:    "waveOutOpen",
			mmresult: mmresult(r),
		}
	}
	return w, nil
}

func waveOutClose(hwo uintptr) error {
	r, _, e := procWaveOutClose.Call(hwo)
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutClose",
			errno: e.(windows.Errno),
		}
	}
	// WAVERR_STILLPLAYING is ignored.
	if mmresult(r) != mmsyserrNoerror && mmresult(r) != waveerrStillplaying {
		return &winmmError{
			fname:    "waveOutClose",
			mmresult: mmresult(r),
		}
	}
	return nil
}

func waveOutPrepareHeader(hwo uintptr, pwh *wavehdr) error {
	r, _, e := procWaveOutPrepareHeader.Call(hwo, uintptr(unsafe.Pointer(pwh)), unsafe.Sizeof(wavehdr{}))
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutPrepareHeader",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return &winmmError{
			fname:    "waveOutPrepareHeader",
			mmresult: mmresult(r),
		}
	}
	return nil
}

func waveOutWrite(hwo uintptr, pwh *wavehdr) error {
	r, _, e := procWaveOutWrite.Call(hwo, uintptr(unsafe.Pointer(pwh)), unsafe.Sizeof(wavehdr{}))
	if e.(windows.Errno) != 0 {
		return &winmmError{
			fname: "waveOutWrite",
			errno: e.(windows.Errno),
		}
	}
	if mmresult(r) != mmsyserrNoerror {
		return &winmmError{
			fname:    "waveOutWrite",
			mmresult: mmresult(r),
		}
	}
	return nil
}
