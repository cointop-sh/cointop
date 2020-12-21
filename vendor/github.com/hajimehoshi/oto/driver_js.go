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

// +build js

package oto

import (
	"errors"

	"github.com/gopherjs/gopherwasm/js"
)

type driver struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int
	nextPos         float64
	tmp             []byte
	bufferSize      int
	context         js.Value
	lastTime        float64
	lastAudioTime   float64
	ready           bool
}

const audioBufferSamples = 3200

func newDriver(sampleRate, channelNum, bitDepthInBytes, bufferSize int) (*driver, error) {
	class := js.Global().Get("AudioContext")
	if class == js.Undefined() {
		class = js.Global().Get("webkitAudioContext")
	}
	if class == js.Undefined() {
		return nil, errors.New("oto: audio couldn't be initialized")
	}
	p := &driver{
		sampleRate:      sampleRate,
		channelNum:      channelNum,
		bitDepthInBytes: bitDepthInBytes,
		context:         class.New(),
		bufferSize:      max(bufferSize, audioBufferSamples*channelNum*bitDepthInBytes),
	}

	setCallback := func(event string) {
		var f js.Callback
		f = js.NewCallback(func(arguments []js.Value) {
			if !p.ready {
				p.context.Call("resume")
				p.ready = true
			}
			js.Global().Get("document").Call("removeEventListener", event, f)
		})
		js.Global().Get("document").Call("addEventListener", event, f)
	}

	// Browsers require user interaction to start the audio.
	// https://developers.google.com/web/updates/2017/09/autoplay-policy-changes#webaudio
	setCallback("touchend")
	setCallback("keyup")
	setCallback("mouseup")
	return p, nil
}

func toLR(data []byte) ([]float32, []float32) {
	const max = 1 << 15

	l := make([]float32, len(data)/4)
	r := make([]float32, len(data)/4)
	for i := 0; i < len(data)/4; i++ {
		l[i] = float32(int16(data[4*i])|int16(data[4*i+1])<<8) / max
		r[i] = float32(int16(data[4*i+2])|int16(data[4*i+3])<<8) / max
	}
	return l, r
}

func nowInSeconds() float64 {
	return js.Global().Get("performance").Call("now").Float() / 1000.0
}

func (p *driver) TryWrite(data []byte) (int, error) {
	if !p.ready {
		return 0, nil
	}

	n := min(len(data), max(0, p.bufferSize-len(p.tmp)))
	p.tmp = append(p.tmp, data[:n]...)

	c := p.context.Get("currentTime").Float()
	now := nowInSeconds()

	if p.lastTime != 0 && p.lastAudioTime != 0 && p.lastAudioTime >= c && p.lastTime != now {
		// Unfortunately, currentTime might not be precise enough on some devices
		// (e.g. Android Chrome). Adjust the audio time with OS clock.
		c = p.lastAudioTime + now - p.lastTime
	}

	p.lastAudioTime = c
	p.lastTime = now

	if p.nextPos < c {
		p.nextPos = c
	}

	// It's too early to enqueue a buffer.
	// Highly likely, there are two playing buffers now.
	if c+float64(p.bufferSize/p.bitDepthInBytes/p.channelNum)/float64(p.sampleRate) < p.nextPos {
		return n, nil
	}

	le := audioBufferSamples * p.bitDepthInBytes * p.channelNum
	if len(p.tmp) < le {
		return n, nil
	}

	buf := p.context.Call("createBuffer", p.channelNum, audioBufferSamples, p.sampleRate)
	l, r := toLR(p.tmp[:le])
	tl := js.TypedArrayOf(l)
	tr := js.TypedArrayOf(r)
	if buf.Get("copyToChannel") != js.Undefined() {
		buf.Call("copyToChannel", tl, 0, 0)
		buf.Call("copyToChannel", tr, 1, 0)
	} else {
		// copyToChannel is not defined on Safari 11
		buf.Call("getChannelData", 0).Call("set", tl)
		buf.Call("getChannelData", 1).Call("set", tr)
	}
	tl.Release()
	tr.Release()

	s := p.context.Call("createBufferSource")
	s.Set("buffer", buf)
	s.Call("connect", p.context.Get("destination"))
	s.Call("start", p.nextPos)
	p.nextPos += buf.Get("duration").Float()

	p.tmp = p.tmp[le:]
	return n, nil
}

func (p *driver) Close() error {
	return nil
}
