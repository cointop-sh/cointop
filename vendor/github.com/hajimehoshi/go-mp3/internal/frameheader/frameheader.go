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

package frameheader

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/go-mp3/internal/consts"
)

// A mepg1FrameHeader is MPEG1 Layer 1-3 frame header
type FrameHeader uint32

// ID returns this header's ID stored in position 20,19
func (f FrameHeader) ID() consts.Version {
	return consts.Version((f & 0x00180000) >> 19)
}

// Layer returns the mpeg layer of this frame stored in position 18,17
func (f FrameHeader) Layer() consts.Layer {
	return consts.Layer((f & 0x00060000) >> 17)
}

// ProtectionBit returns the protection bit stored in position 16
func (f FrameHeader) ProtectionBit() int {
	return int(f&0x00010000) >> 16
}

// BirateIndex returns the bitrate index stored in position 15,12
func (f FrameHeader) BitrateIndex() int {
	return int(f&0x0000f000) >> 12
}

// SamplingFrequency returns the SamplingFrequency in Hz stored in position 11,10
func (f FrameHeader) SamplingFrequency() consts.SamplingFrequency {
	return consts.SamplingFrequency(int(f&0x00000c00) >> 10)
}

// PaddingBit returns the padding bit stored in position 9
func (f FrameHeader) PaddingBit() int {
	return int(f&0x00000200) >> 9
}

// PrivateBit returns the private bit stored in position 8 - this bit may be used to store arbitrary data to be used
// by an application
func (f FrameHeader) PrivateBit() int {
	return int(f&0x00000100) >> 8
}

// Mode returns the channel mode, stored in position 7,6
func (f FrameHeader) Mode() consts.Mode {
	return consts.Mode((f & 0x000000c0) >> 6)
}

// modeExtension returns the mode_extension - for use with Joint Stereo - stored in position 4,5
func (f FrameHeader) modeExtension() int {
	return int(f&0x00000030) >> 4
}

// UseMSStereo returns a boolean value indicating whether the frame uses middle/side stereo.
func (f FrameHeader) UseMSStereo() bool {
	if f.Mode() != consts.ModeJointStereo {
		return false
	}
	return f.modeExtension()&0x2 != 0
}

// UseIntensityStereo returns a boolean value indicating whether the frame uses intensity stereo.
func (f FrameHeader) UseIntensityStereo() bool {
	if f.Mode() != consts.ModeJointStereo {
		return false
	}
	return f.modeExtension()&0x1 != 0
}

// Copyright returns whether or not this recording is copywritten - stored in position 3
func (f FrameHeader) Copyright() int {
	return int(f&0x00000008) >> 3
}

// OriginalOrCopy returns whether or not this is an Original recording or a copy of one - stored in position 2
func (f FrameHeader) OriginalOrCopy() int {
	return int(f&0x00000004) >> 2
}

// Emphasis returns emphasis - the emphasis indication is here to tell the decoder that the file must be de-emphasized - stored in position 0,1
func (f FrameHeader) Emphasis() int {
	return int(f&0x00000003) >> 0
}

// IsValid returns a boolean value indicating whether the header is valid or not.
func (f FrameHeader) IsValid() bool {
	const sync = 0xffe00000
	if (f & sync) != sync {
		return false
	}
	if f.ID() == consts.VersionReserved {
		return false
	}
	if f.BitrateIndex() == 15 {
		return false
	}
	if f.SamplingFrequency() == consts.SamplingFrequencyReserved {
		return false
	}
	if f.Layer() == consts.LayerReserved {
		return false
	}
	if f.Emphasis() == 2 {
		return false
	}
	return true
}

func bitrate(layer consts.Layer, index int) int {
	switch layer {
	case consts.Layer1:
		return []int{
			0, 32000, 64000, 96000, 128000, 160000, 192000, 224000,
			256000, 288000, 320000, 352000, 384000, 416000, 448000}[index]
	case consts.Layer2:
		return []int{
			0, 32000, 48000, 56000, 64000, 80000, 96000, 112000,
			128000, 160000, 192000, 224000, 256000, 320000, 384000}[index]
	case consts.Layer3:
		return []int{
			0, 32000, 40000, 48000, 56000, 64000, 80000, 96000,
			112000, 128000, 160000, 192000, 224000, 256000, 320000}[index]
	}
	panic("not reached")
}

func (f FrameHeader) FrameSize() int {
	return (144*bitrate(f.Layer(), f.BitrateIndex()))/
		f.SamplingFrequency().Int() +
		int(f.PaddingBit())
}

func (f FrameHeader) NumberOfChannels() int {
	if f.Mode() == consts.ModeSingleChannel {
		return 1
	}
	return 2
}

type FullReader interface {
	ReadFull([]byte) (int, error)
}

func Read(source FullReader, position int64) (h FrameHeader, startPosition int64, err error) {
	buf := make([]byte, 4)
	if n, err := source.ReadFull(buf); n < 4 {
		if err == io.EOF {
			if n == 0 {
				// Expected EOF
				return 0, 0, io.EOF
			}
			return 0, 0, &consts.UnexpectedEOF{"readHeader (1)"}
		}
		return 0, 0, err
	}

	b1 := uint32(buf[0])
	b2 := uint32(buf[1])
	b3 := uint32(buf[2])
	b4 := uint32(buf[3])
	header := FrameHeader((b1 << 24) | (b2 << 16) | (b3 << 8) | (b4 << 0))
	for !header.IsValid() {
		b1 = b2
		b2 = b3
		b3 = b4

		buf := make([]byte, 1)
		if _, err := source.ReadFull(buf); err != nil {
			if err == io.EOF {
				return 0, 0, &consts.UnexpectedEOF{"readHeader (2)"}
			}
			return 0, 0, err
		}
		b4 = uint32(buf[0])
		header = FrameHeader((b1 << 24) | (b2 << 16) | (b3 << 8) | (b4 << 0))
		position++
	}

	// If we get here we've found the sync word, and can decode the header
	// which is in the low 20 bits of the 32-bit sync+header word.

	if header.BitrateIndex() == 0 {
		return 0, 0, fmt.Errorf("mp3: free bitrate format is not supported. Header word is 0x%08x at position %d",
			header, position)
	}
	return header, position, nil
}
