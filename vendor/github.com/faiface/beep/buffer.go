package beep

import (
	"fmt"
	"math"
	"time"
)

// SampleRate is the number of samples per second.
type SampleRate int

// D returns the duration of n samples.
func (sr SampleRate) D(n int) time.Duration {
	return time.Second * time.Duration(n) / time.Duration(sr)
}

// N returns the number of samples that last for d duration.
func (sr SampleRate) N(d time.Duration) int {
	return int(d * time.Duration(sr) / time.Second)
}

// Format is the format of a Buffer or another audio source.
type Format struct {
	// SampleRate is the number of samples per second.
	SampleRate SampleRate

	// NumChannels is the number of channels. The value of 1 is mono, the value of 2 is stereo.
	// The samples should always be interleaved.
	NumChannels int

	// Precision is the number of bytes used to encode a single sample. Only values up to 6 work
	// well, higher values loose precision due to floating point numbers.
	Precision int
}

// Width returns the number of bytes per one frame (samples in all channels).
//
// This is equal to f.NumChannels * f.Precision.
func (f Format) Width() int {
	return f.NumChannels * f.Precision
}

// EncodeSigned encodes a single sample in f.Width() bytes to p in signed format.
func (f Format) EncodeSigned(p []byte, sample [2]float64) (n int) {
	return f.encode(true, p, sample)
}

// EncodeUnsigned encodes a single sample in f.Width() bytes to p in unsigned format.
func (f Format) EncodeUnsigned(p []byte, sample [2]float64) (n int) {
	return f.encode(false, p, sample)
}

// DecodeSigned decodes a single sample encoded in f.Width() bytes from p in signed format.
func (f Format) DecodeSigned(p []byte) (sample [2]float64, n int) {
	return f.decode(true, p)
}

// DecodeUnsigned decodes a single sample encoded in f.Width() bytes from p in unsigned format.
func (f Format) DecodeUnsigned(p []byte) (sample [2]float64, n int) {
	return f.decode(false, p)
}

func (f Format) encode(signed bool, p []byte, sample [2]float64) (n int) {
	switch {
	case f.NumChannels == 1:
		x := norm((sample[0] + sample[1]) / 2)
		p = p[encodeFloat(signed, f.Precision, p, x):]
	case f.NumChannels >= 2:
		for c := range sample {
			x := norm(sample[c])
			p = p[encodeFloat(signed, f.Precision, p, x):]
		}
		for c := len(sample); c < f.NumChannels; c++ {
			p = p[encodeFloat(signed, f.Precision, p, 0):]
		}
	default:
		panic(fmt.Errorf("format: encode: invalid number of channels: %d", f.NumChannels))
	}
	return f.Width()
}

func (f Format) decode(signed bool, p []byte) (sample [2]float64, n int) {
	switch {
	case f.NumChannels == 1:
		x, _ := decodeFloat(signed, f.Precision, p)
		return [2]float64{x, x}, f.Width()
	case f.NumChannels >= 2:
		for c := range sample {
			x, n := decodeFloat(signed, f.Precision, p)
			sample[c] = x
			p = p[n:]
		}
		for c := len(sample); c < f.NumChannels; c++ {
			_, n := decodeFloat(signed, f.Precision, p)
			p = p[n:]
		}
		return sample, f.Width()
	default:
		panic(fmt.Errorf("format: decode: invalid number of channels: %d", f.NumChannels))
	}
}

func encodeFloat(signed bool, precision int, p []byte, x float64) (n int) {
	var xUint64 uint64
	if signed {
		xUint64 = floatToSigned(precision, x)
	} else {
		xUint64 = floatToUnsigned(precision, x)
	}
	for i := 0; i < precision; i++ {
		p[i] = byte(xUint64)
		xUint64 >>= 8
	}
	return precision
}

func decodeFloat(signed bool, precision int, p []byte) (x float64, n int) {
	var xUint64 uint64
	for i := precision - 1; i >= 0; i-- {
		xUint64 <<= 8
		xUint64 += uint64(p[i])
	}
	if signed {
		return signedToFloat(precision, xUint64), precision
	}
	return unsignedToFloat(precision, xUint64), precision
}

func floatToSigned(precision int, x float64) uint64 {
	if x < 0 {
		compl := uint64(-x * (math.Exp2(float64(precision)*8-1) - 1))
		return uint64(1<<uint(precision*8)) - compl
	}
	return uint64(x * (math.Exp2(float64(precision)*8-1) - 1))
}

func floatToUnsigned(precision int, x float64) uint64 {
	return uint64((x + 1) / 2 * (math.Exp2(float64(precision)*8) - 1))
}

func signedToFloat(precision int, xUint64 uint64) float64 {
	if xUint64 >= 1<<uint(precision*8-1) {
		compl := 1<<uint(precision*8) - xUint64
		return -float64(int64(compl)) / (math.Exp2(float64(precision)*8-1) - 1)
	}
	return float64(int64(xUint64)) / (math.Exp2(float64(precision)*8-1) - 1)
}

func unsignedToFloat(precision int, xUint64 uint64) float64 {
	return float64(xUint64)/(math.Exp2(float64(precision)*8)-1)*2 - 1
}

func norm(x float64) float64 {
	if x < -1 {
		return -1
	}
	if x > +1 {
		return +1
	}
	return x
}

// Buffer is a storage for audio data. You can think of it as a bytes.Buffer for audio samples.
type Buffer struct {
	f    Format
	data []byte
	tmp  []byte
}

// NewBuffer creates a new empty Buffer which stores samples in the provided format.
func NewBuffer(f Format) *Buffer {
	return &Buffer{f: f, tmp: make([]byte, f.Width())}
}

// Format returns the format of the Buffer.
func (b *Buffer) Format() Format {
	return b.f
}

// Len returns the number of samples currently in the Buffer.
func (b *Buffer) Len() int {
	return len(b.data) / b.f.Width()
}

// Pop removes n samples from the beginning of the Buffer.
//
// Existing Streamers are not affected.
func (b *Buffer) Pop(n int) {
	b.data = b.data[n*b.f.Width():]
}

// Append adds all audio data from the given Streamer to the end of the Buffer.
//
// The Streamer will be drained when this method finishes.
func (b *Buffer) Append(s Streamer) {
	var samples [512][2]float64
	for {
		n, ok := s.Stream(samples[:])
		if !ok {
			break
		}
		for _, sample := range samples[:n] {
			b.f.EncodeSigned(b.tmp, sample)
			b.data = append(b.data, b.tmp...)
		}
	}
}

// Streamer returns a StreamSeeker which streams samples in the given interval (including from,
// excluding to). If from<0 or to>b.Len() or to<from, this method panics.
//
// When using multiple goroutines, synchronization of Streamers with the Buffer is not required,
// as Buffer is persistent (but efficient and garbage collected).
func (b *Buffer) Streamer(from, to int) StreamSeeker {
	return &bufferStreamer{
		f:    b.f,
		data: b.data[from*b.f.Width() : to*b.f.Width()],
		pos:  0,
	}
}

type bufferStreamer struct {
	f    Format
	data []byte
	pos  int
}

func (bs *bufferStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if bs.pos >= len(bs.data) {
		return 0, false
	}
	for i := range samples {
		if bs.pos >= len(bs.data) {
			break
		}
		sample, advance := bs.f.DecodeSigned(bs.data[bs.pos:])
		samples[i] = sample
		bs.pos += advance
		n++
	}
	return n, true
}

func (bs *bufferStreamer) Err() error {
	return nil
}

func (bs *bufferStreamer) Len() int {
	return len(bs.data) / bs.f.Width()
}

func (bs *bufferStreamer) Position() int {
	return bs.pos / bs.f.Width()
}

func (bs *bufferStreamer) Seek(p int) error {
	if p < 0 || bs.Len() < p {
		return fmt.Errorf("buffer: seek position %v out of range [%v, %v]", p, 0, bs.Len())
	}
	bs.pos = p * bs.f.Width()
	return nil
}
