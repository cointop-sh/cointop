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

// Package oto offers io.Writer to play sound on multiple platforms.
package oto

import (
	"io"
	"runtime"
)

// Player is a PCM (pulse-code modulation) audio player.
// Player implements io.WriteCloser.
// Use Write method to play samples.
type Player struct {
	context *Context
	r       *io.PipeReader
	w       *io.PipeWriter
}

func newPlayer(context *Context) *Player {
	r, w := io.Pipe()
	p := &Player{
		context: context,
		r:       r,
		w:       w,
	}
	runtime.SetFinalizer(p, (*Player).Close)
	return p
}

// Write writes PCM samples to the Player.
//
// The format is as follows:
//   [data]      = [sample 1] [sample 2] [sample 3] ...
//   [sample *]  = [channel 1] ...
//   [channel *] = [byte 1] [byte 2] ...
// Byte ordering is little endian.
//
// The data is first put into the Player's buffer. Once the buffer is full, Player starts playing
// the data and empties the buffer.
//
// If the supplied data doesn't fit into the Player's buffer, Write block until a sufficient amount
// of data has been played (or at least started playing) and the remaining unplayed data fits into
// the buffer.
//
// Note, that the Player won't start playing anything until the buffer is full.
func (p *Player) Write(buf []byte) (int, error) {
	return p.w.Write(buf)
}

// Close closes the Player and frees any resources associated with it. The Player is no longer
// usable after calling Close.
func (p *Player) Close() error {
	runtime.SetFinalizer(p, nil)

	// Already closed
	if p.context == nil {
		return nil
	}

	// Close the pipe writer before RemoveSource, or Read-ing in the mux takes forever.
	if err := p.w.Close(); err != nil {
		return err
	}

	p.context.mux.RemoveSource(p.r)
	p.context = nil

	// Close the pipe reader after RemoveSource, or ErrClosedPipe happens at Read-ing.
	if err := p.r.Close(); err != nil {
		return err
	}
	return nil
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
