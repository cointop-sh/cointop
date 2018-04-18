// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"sync"

	tm "github.com/nsf/termbox-go"
)

// Bufferer should be implemented by all renderable components.
type Bufferer interface {
	Buffer() Buffer
}

// Close finalizes termui library,
// should be called after successful initialization when termui's functionality isn't required anymore.
func Close() {
	tm.Close()
}

var renderLock sync.Mutex

func termSync() {
	renderLock.Lock()
	tm.Sync()
	termWidth, termHeight = tm.Size()
	renderLock.Unlock()
}

// TermWidth returns the current terminal's width.
func TermWidth() int {
	termSync()
	return termWidth
}

// TermHeight returns the current terminal's height.
func TermHeight() int {
	termSync()
	return termHeight
}

func Clear() {
	tm.Clear(tm.ColorDefault, toTmAttr(ThemeAttr("bg")))
}

func clearArea(r image.Rectangle, bg Attribute) {
	for i := r.Min.X; i < r.Max.X; i++ {
		for j := r.Min.Y; j < r.Max.Y; j++ {
			tm.SetCell(i, j, ' ', tm.ColorDefault, toTmAttr(bg))
		}
	}
}

func ClearArea(r image.Rectangle, bg Attribute) {
	clearArea(r, bg)
	tm.Flush()
}

var renderJobs chan []Bufferer

func Render(bs ...Bufferer) {
	//go func() { renderJobs <- bs }()
	renderJobs <- bs
}
