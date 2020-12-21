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

	"github.com/hajimehoshi/go-mp3/internal/bits"
	"github.com/hajimehoshi/go-mp3/internal/consts"
	"github.com/hajimehoshi/go-mp3/internal/frameheader"
	"github.com/hajimehoshi/go-mp3/internal/huffman"
	"github.com/hajimehoshi/go-mp3/internal/sideinfo"
)

func readHuffman(m *bits.Bits, header frameheader.FrameHeader, sideInfo *sideinfo.SideInfo, mainData *MainData, part_2_start, gr, ch int) error {
	// Check that there is any data to decode. If not, zero the array.
	if sideInfo.Part2_3Length[gr][ch] == 0 {
		for i := 0; i < consts.SamplesPerGr; i++ {
			mainData.Is[gr][ch][i] = 0.0
		}
		return nil
	}

	// Calculate bit_pos_end which is the index of the last bit for this part.
	bit_pos_end := part_2_start + sideInfo.Part2_3Length[gr][ch] - 1
	// Determine region boundaries
	region_1_start := 0
	region_2_start := 0
	if (sideInfo.WinSwitchFlag[gr][ch] == 1) && (sideInfo.BlockType[gr][ch] == 2) {
		region_1_start = 36                  // sfb[9/3]*3=36
		region_2_start = consts.SamplesPerGr // No Region2 for short block case.
	} else {
		sfreq := header.SamplingFrequency()
		l := consts.SfBandIndicesSet[sfreq].L
		i := sideInfo.Region0Count[gr][ch] + 1
		if i < 0 || len(l) <= i {
			// TODO: Better error messages (#3)
			return fmt.Errorf("mp3: readHuffman failed: invalid index i: %d", i)
		}
		region_1_start = l[i]
		j := sideInfo.Region0Count[gr][ch] + sideInfo.Region1Count[gr][ch] + 2
		if j < 0 || len(l) <= j {
			// TODO: Better error messages (#3)
			return fmt.Errorf("mp3: readHuffman failed: invalid index j: %d", j)
		}
		region_2_start = l[j]
	}
	// Read big_values using tables according to region_x_start
	for is_pos := 0; is_pos < sideInfo.BigValues[gr][ch]*2; is_pos++ {
		// #22
		if is_pos >= len(mainData.Is[gr][ch]) {
			return fmt.Errorf("mp3: is_pos was too big: %d", is_pos)
		}
		table_num := 0
		if is_pos < region_1_start {
			table_num = sideInfo.TableSelect[gr][ch][0]
		} else if is_pos < region_2_start {
			table_num = sideInfo.TableSelect[gr][ch][1]
		} else {
			table_num = sideInfo.TableSelect[gr][ch][2]
		}
		// Get next Huffman coded words
		x, y, _, _, err := huffman.Decode(m, table_num)
		if err != nil {
			return err
		}
		// In the big_values area there are two freq lines per Huffman word
		mainData.Is[gr][ch][is_pos] = float32(x)
		is_pos++
		mainData.Is[gr][ch][is_pos] = float32(y)
	}
	// Read small values until is_pos = 576 or we run out of huffman data
	// TODO: Is this comment wrong?
	table_num := sideInfo.Count1TableSelect[gr][ch] + 32
	is_pos := sideInfo.BigValues[gr][ch] * 2
	for is_pos <= 572 && m.BitPos() <= bit_pos_end {
		// Get next Huffman coded words
		x, y, v, w, err := huffman.Decode(m, table_num)
		if err != nil {
			return err
		}
		mainData.Is[gr][ch][is_pos] = float32(v)
		is_pos++
		if is_pos >= consts.SamplesPerGr {
			break
		}
		mainData.Is[gr][ch][is_pos] = float32(w)
		is_pos++
		if is_pos >= consts.SamplesPerGr {
			break
		}
		mainData.Is[gr][ch][is_pos] = float32(x)
		is_pos++
		if is_pos >= consts.SamplesPerGr {
			break
		}
		mainData.Is[gr][ch][is_pos] = float32(y)
		is_pos++
	}
	// Check that we didn't read past the end of this section
	if m.BitPos() > (bit_pos_end + 1) {
		// Remove last words read
		is_pos -= 4
	}

	// Setup count1 which is the index of the first sample in the rzero reg.
	sideInfo.Count1[gr][ch] = is_pos

	// Zero out the last part if necessary
	for is_pos < consts.SamplesPerGr {
		mainData.Is[gr][ch][is_pos] = 0.0
		is_pos++
	}
	// Set the bitpos to point to the next part to read
	m.SetPos(bit_pos_end + 1)
	return nil
}
