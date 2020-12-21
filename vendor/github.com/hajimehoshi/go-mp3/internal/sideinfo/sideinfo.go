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

package sideinfo

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/go-mp3/internal/bits"
	"github.com/hajimehoshi/go-mp3/internal/consts"
	"github.com/hajimehoshi/go-mp3/internal/frameheader"
)

type FullReader interface {
	ReadFull([]byte) (int, error)
}

// A SideInfo is MPEG1 Layer 3 Side Information.
// [2][2] means [gr][ch].
type SideInfo struct {
	MainDataBegin    int       // 9 bits
	PrivateBits      int       // 3 bits in mono, 5 in stereo
	Scfsi            [2][4]int // 1 bit
	Part2_3Length    [2][2]int // 12 bits
	BigValues        [2][2]int // 9 bits
	GlobalGain       [2][2]int // 8 bits
	ScalefacCompress [2][2]int // 4 bits
	WinSwitchFlag    [2][2]int // 1 bit

	BlockType      [2][2]int    // 2 bits
	MixedBlockFlag [2][2]int    // 1 bit
	TableSelect    [2][2][3]int // 5 bits
	SubblockGain   [2][2][3]int // 3 bits

	Region0Count [2][2]int // 4 bits
	Region1Count [2][2]int // 3 bits

	Preflag           [2][2]int // 1 bit
	ScalefacScale     [2][2]int // 1 bit
	Count1TableSelect [2][2]int // 1 bit
	Count1            [2][2]int // Not in file, calc by huffman decoder
}

func Read(source FullReader, header frameheader.FrameHeader) (*SideInfo, error) {
	nch := header.NumberOfChannels()
	// Calculate header audio data size
	framesize := header.FrameSize()
	if framesize > 2000 {
		return nil, fmt.Errorf("mp3: framesize = %d\n", framesize)
	}
	// Sideinfo is 17 bytes for one channel and 32 bytes for two
	sideinfo_size := 32
	if nch == 1 {
		sideinfo_size = 17
	}
	// Main data size is the rest of the frame,including ancillary data
	main_data_size := framesize - sideinfo_size - 4 // sync+header
	// CRC is 2 bytes
	if header.ProtectionBit() == 0 {
		main_data_size -= 2
	}
	// Read sideinfo from bitstream into buffer used by Bits()
	buf := make([]byte, sideinfo_size)
	n, err := source.ReadFull(buf)
	if n < sideinfo_size {
		if err == io.EOF {
			return nil, &consts.UnexpectedEOF{"sideinfo.Read"}
		}
		return nil, fmt.Errorf("mp3: couldn't read sideinfo %d bytes: %v", sideinfo_size, err)
	}
	s := bits.New(buf)

	// Parse audio data
	// Pointer to where we should start reading main data
	si := &SideInfo{}
	si.MainDataBegin = s.Bits(9)
	// Get private bits. Not used for anything.
	if header.Mode() == consts.ModeSingleChannel {
		si.PrivateBits = s.Bits(5)
	} else {
		si.PrivateBits = s.Bits(3)
	}
	// Get scale factor selection information
	for ch := 0; ch < nch; ch++ {
		for scfsi_band := 0; scfsi_band < 4; scfsi_band++ {
			si.Scfsi[ch][scfsi_band] = s.Bits(1)
		}
	}
	// Get the rest of the side information
	for gr := 0; gr < 2; gr++ {
		for ch := 0; ch < nch; ch++ {
			si.Part2_3Length[gr][ch] = s.Bits(12)
			si.BigValues[gr][ch] = s.Bits(9)
			si.GlobalGain[gr][ch] = s.Bits(8)
			si.ScalefacCompress[gr][ch] = s.Bits(4)
			si.WinSwitchFlag[gr][ch] = s.Bits(1)
			if si.WinSwitchFlag[gr][ch] == 1 {
				si.BlockType[gr][ch] = s.Bits(2)
				si.MixedBlockFlag[gr][ch] = s.Bits(1)
				for region := 0; region < 2; region++ {
					si.TableSelect[gr][ch][region] = s.Bits(5)
				}
				for window := 0; window < 3; window++ {
					si.SubblockGain[gr][ch][window] = s.Bits(3)
				}

				// TODO: This is not listed on the spec. Is this correct??
				if si.BlockType[gr][ch] == 2 && si.MixedBlockFlag[gr][ch] == 0 {
					si.Region0Count[gr][ch] = 8 // Implicit
				} else {
					si.Region0Count[gr][ch] = 7 // Implicit
				}
				// The standard is wrong on this!!!
				// Implicit
				si.Region1Count[gr][ch] = 20 - si.Region0Count[gr][ch]
			} else {
				for region := 0; region < 3; region++ {
					si.TableSelect[gr][ch][region] = s.Bits(5)
				}
				si.Region0Count[gr][ch] = s.Bits(4)
				si.Region1Count[gr][ch] = s.Bits(3)
				si.BlockType[gr][ch] = 0 // Implicit
			}
			si.Preflag[gr][ch] = s.Bits(1)
			si.ScalefacScale[gr][ch] = s.Bits(1)
			si.Count1TableSelect[gr][ch] = s.Bits(1)
		}
	}
	return si, nil
}
