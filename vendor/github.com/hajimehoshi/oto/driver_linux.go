// Copyright 2017 The Oto Authors
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
// +build !android
// +build !ios

package oto

/*
#cgo LDFLAGS: -lasound
#include <alsa/asoundlib.h>

static void check(int *err, int newErr) {
  if (*err) {
    return;
  }
  *err = newErr;
}

static int ALSA_hw_params(
    snd_pcm_t          *pcm,
    unsigned           sampleRate,
    unsigned           numChans,
    snd_pcm_format_t   format,
    snd_pcm_uframes_t* buffer_size,
    snd_pcm_uframes_t* period_size) {
  snd_pcm_hw_params_t* params = NULL;
  int err = 0;
  snd_pcm_hw_params_alloca(&params);
  check(&err, snd_pcm_hw_params_any(pcm, params));

  check(&err, snd_pcm_hw_params_set_access(pcm, params, SND_PCM_ACCESS_RW_INTERLEAVED));
  check(&err, snd_pcm_hw_params_set_format(pcm, params, format));
  check(&err, snd_pcm_hw_params_set_channels(pcm, params, numChans));
  check(&err, snd_pcm_hw_params_set_rate_resample(pcm, params, 1));
  check(&err, snd_pcm_hw_params_set_rate_near(pcm, params, &sampleRate, NULL));
  check(&err, snd_pcm_hw_params_set_buffer_size_near(pcm, params, buffer_size));
  check(&err, snd_pcm_hw_params_set_period_size_near(pcm, params, period_size, NULL));

  check(&err, snd_pcm_hw_params(pcm, params));

  return err;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type driver struct {
	handle          *C.snd_pcm_t
	buf             []byte
	bufSamples      int
	numChans        int
	bitDepthInBytes int
}

func alsaError(err C.int) error {
	return fmt.Errorf("oto: ALSA error: %s", C.GoString(C.snd_strerror(err)))
}

func newDriver(sampleRate, numChans, bitDepthInBytes, bufferSizeInBytes int) (*driver, error) {
	p := &driver{
		numChans:        numChans,
		bitDepthInBytes: bitDepthInBytes,
	}

	// open a default ALSA audio device for blocking stream playback
	if errCode := C.snd_pcm_open(&p.handle, C.CString("default"), C.SND_PCM_STREAM_PLAYBACK, 0); errCode < 0 {
		return nil, alsaError(errCode)
	}

	// bufferSize is the total size of the main circular buffer fullness of this buffer
	// oscilates somewhere between bufferSize and bufferSize-periodSize
	bufferSize := C.snd_pcm_uframes_t(bufferSizeInBytes / (numChans * bitDepthInBytes))
	// periodSize is the number of samples that will be taken from the main circular
	// buffer at once, we leave this value to bufferSize, because ALSA will change that
	// to the maximum viable number, obviously lower than bufferSize
	periodSize := bufferSize

	// choose the correct sample format according to bitDepthInBytes
	var format C.snd_pcm_format_t
	switch bitDepthInBytes {
	case 1:
		format = C.SND_PCM_FORMAT_S8
	case 2:
		format = C.SND_PCM_FORMAT_S16_LE
	default:
		panic(fmt.Errorf("oto: bitDepthInBytes must be 1 or 2, got %d", bitDepthInBytes))
	}

	// set the device hardware parameters according to sampleRate, numChans, format, bufferSize
	// and periodSize
	//
	// bufferSize and periodSize are passed as pointers, because they may be changed according
	// to the wisdom of ALSA
	//
	// ALSA will try too keep them as close to what was requested as possible
	if errCode := C.ALSA_hw_params(p.handle, C.uint(sampleRate), C.uint(numChans), format, &bufferSize, &periodSize); errCode < 0 {
		p.Close()
		return nil, alsaError(errCode)
	}

	// allocate the buffer of the size of the period, use the periodSize that we've got back
	// from ALSA after it's wise decision
	p.bufSamples = int(periodSize)
	p.buf = []byte{}

	return p, nil
}

func (p *driver) TryWrite(data []byte) (n int, err error) {
	bufSize := p.bufSamples * p.numChans * p.bitDepthInBytes
	for len(data) > 0 {
		toWrite := min(len(data), max(0, bufSize-len(p.buf)))
		p.buf = append(p.buf, data[:toWrite]...)
		data = data[toWrite:]
		n += toWrite

		// our buffer is not full and we've used up all the data, we'll keep them and finish
		if len(p.buf) < bufSize {
			break
		}

		// write samples to the main circular buffer
		wrote := C.snd_pcm_writei(p.handle, unsafe.Pointer(&p.buf[0]), C.snd_pcm_uframes_t(p.bufSamples))
		if wrote == -C.EPIPE {
			// Underrun!
			if errCode := C.snd_pcm_prepare(p.handle); errCode < 0 {
				return 0, alsaError(errCode)
			}
			continue
		}
		if wrote < 0 {
			// an error occured while writing samples
			return 0, alsaError(C.int(wrote))
		}
		p.buf = p.buf[int(wrote)*p.numChans*p.bitDepthInBytes:]
	}
	return n, nil
}

func (p *driver) Close() error {
	// drop the remaining unprocessed samples in the main circular buffer
	if errCode := C.snd_pcm_drop(p.handle); errCode < 0 {
		return alsaError(errCode)
	}
	if errCode := C.snd_pcm_close(p.handle); errCode < 0 {
		return alsaError(errCode)
	}
	return nil
}
