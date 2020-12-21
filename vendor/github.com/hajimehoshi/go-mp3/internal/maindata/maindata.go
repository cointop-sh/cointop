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

package maindata

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/go-mp3/internal/bits"
	"github.com/hajimehoshi/go-mp3/internal/consts"
	"github.com/hajimehoshi/go-mp3/internal/frameheader"
	"github.com/hajimehoshi/go-mp3/internal/sideinfo"
)

type FullReader interface {
	ReadFull([]byte) (int, error)
}

// A MainData is MPEG1 Layer 3 Main Data.
type MainData struct {
	ScalefacL [2][2][22]int      // 0-4 bits
	ScalefacS [2][2][13][3]int   // 0-4 bits
	Is        [2][2][576]float32 // Huffman coded freq. lines
}

var scalefacSizes = [16][2]int{
	{0, 0}, {0, 1}, {0, 2}, {0, 3}, {3, 0}, {1, 1}, {1, 2}, {1, 3},
	{2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}, {4, 2}, {4, 3},
}

func Read(source FullReader, prev *bits.Bits, header frameheader.FrameHeader, sideInfo *sideinfo.SideInfo) (*MainData, *bits.Bits, error) {
	nch := header.NumberOfChannels()
	// Calculate header audio data size
	framesize := header.FrameSize()
	if framesize > 2000 {
		return nil, nil, fmt.Errorf("mp3: framesize = %d", framesize)
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
	// Assemble main data buffer with data from this frame and the previous
	// two frames. main_data_begin indicates how many bytes from previous
	// frames that should be used. This buffer is later accessed by the
	// Bits function in the same way as the side info is.
	m, err := read(source, prev, main_data_size, sideInfo.MainDataBegin)
	if err != nil {
		// This could be due to not enough data in reservoir
		return nil, nil, err
	}
	md := &MainData{}
	for gr := 0; gr < 2; gr++ {
		for ch := 0; ch < nch; ch++ {
			part_2_start := m.BitPos()
			// Number of bits in the bitstream for the bands
			slen1 := scalefacSizes[sideInfo.ScalefacCompress[gr][ch]][0]
			slen2 := scalefacSizes[sideInfo.ScalefacCompress[gr][ch]][1]
			if sideInfo.WinSwitchFlag[gr][ch] == 1 && sideInfo.BlockType[gr][ch] == 2 {
				if sideInfo.MixedBlockFlag[gr][ch] != 0 {
					for sfb := 0; sfb < 8; sfb++ {
						md.ScalefacL[gr][ch][sfb] = m.Bits(slen1)
					}
					for sfb := 3; sfb < 12; sfb++ {
						//slen1 for band 3-5,slen2 for 6-11
						nbits := slen2
						if sfb < 6 {
							nbits = slen1
						}
						for win := 0; win < 3; win++ {
							md.ScalefacS[gr][ch][sfb][win] = m.Bits(nbits)
						}
					}
				} else {
					for sfb := 0; sfb < 12; sfb++ {
						//slen1 for band 3-5,slen2 for 6-11
						nbits := slen2
						if sfb < 6 {
							nbits = slen1
						}
						for win := 0; win < 3; win++ {
							md.ScalefacS[gr][ch][sfb][win] = m.Bits(nbits)
						}
					}
				}
			} else {
				// Scale factor bands 0-5
				if sideInfo.Scfsi[ch][0] == 0 || gr == 0 {
					for sfb := 0; sfb < 6; sfb++ {
						md.ScalefacL[gr][ch][sfb] = m.Bits(slen1)
					}
				} else if sideInfo.Scfsi[ch][0] == 1 && gr == 1 {
					// Copy scalefactors from granule 0 to granule 1
					// TODO: This is not listed on the spec.
					for sfb := 0; sfb < 6; sfb++ {
						md.ScalefacL[1][ch][sfb] = md.ScalefacL[0][ch][sfb]
					}
				}
				// Scale factor bands 6-10
				if sideInfo.Scfsi[ch][1] == 0 || gr == 0 {
					for sfb := 6; sfb < 11; sfb++ {
						md.ScalefacL[gr][ch][sfb] = m.Bits(slen1)
					}
				} else if sideInfo.Scfsi[ch][1] == 1 && gr == 1 {
					// Copy scalefactors from granule 0 to granule 1
					for sfb := 6; sfb < 11; sfb++ {
						md.ScalefacL[1][ch][sfb] = md.ScalefacL[0][ch][sfb]
					}
				}
				// Scale factor bands 11-15
				if sideInfo.Scfsi[ch][2] == 0 || gr == 0 {
					for sfb := 11; sfb < 16; sfb++ {
						md.ScalefacL[gr][ch][sfb] = m.Bits(slen2)
					}
				} else if sideInfo.Scfsi[ch][2] == 1 && gr == 1 {
					// Copy scalefactors from granule 0 to granule 1
					for sfb := 11; sfb < 16; sfb++ {
						md.ScalefacL[1][ch][sfb] = md.ScalefacL[0][ch][sfb]
					}
				}
				// Scale factor bands 16-20
				if sideInfo.Scfsi[ch][3] == 0 || gr == 0 {
					for sfb := 16; sfb < 21; sfb++ {
						md.ScalefacL[gr][ch][sfb] = m.Bits(slen2)
					}
				} else if sideInfo.Scfsi[ch][3] == 1 && gr == 1 {
					// Copy scalefactors from granule 0 to granule 1
					for sfb := 16; sfb < 21; sfb++ {
						md.ScalefacL[1][ch][sfb] = md.ScalefacL[0][ch][sfb]
					}
				}
			}
			// Read Huffman coded data. Skip stuffing bits.
			if err := readHuffman(m, header, sideInfo, md, part_2_start, gr, ch); err != nil {
				return nil, nil, err
			}
		}
	}
	// The ancillary data is stored here,but we ignore it.
	return md, m, nil
}

func read(source FullReader, prev *bits.Bits, size int, offset int) (*bits.Bits, error) {
	if size > 1500 {
		return nil, fmt.Errorf("mp3: size = %d", size)
	}
	// Check that there's data available from previous frames if needed
	if prev != nil && offset > prev.LenInBytes() {
		// No, there is not, so we skip decoding this frame, but we have to
		// read the main_data bits from the bitstream in case they are needed
		// for decoding the next frame.
		buf := make([]byte, size)
		if n, err := source.ReadFull(buf); n < size {
			if err == io.EOF {
				return nil, &consts.UnexpectedEOF{"maindata.Read (1)"}
			}
			return nil, err
		}
		// TODO: Define a special error and enable to continue the next frame.
		return bits.Append(prev, buf), nil
	}
	// Copy data from previous frames
	vec := []byte{}
	if prev != nil {
		vec = prev.Tail(offset)
	}
	// Read the main_data from file
	buf := make([]byte, size)
	if n, err := source.ReadFull(buf); n < size {
		if err == io.EOF {
			return nil, &consts.UnexpectedEOF{"maindata.Read (2)"}
		}
		return nil, err
	}
	return bits.New(append(vec, buf...)), nil
}
