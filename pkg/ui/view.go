package ui

import (
	"fmt"

	"github.com/cointop-sh/cointop/pkg/gocui"
	"github.com/cointop-sh/cointop/pkg/termbox"
)

// IView is the view interface
type IView interface {
	Backing() *gocui.View
	SetBacking(gocuiView *gocui.View)
	Name() string
}

// View is a view sruct
type View struct {
	backing *gocui.View
	name    string
}

// NewView creates a new view
func NewView(name string) *View {
	return &View{
		name: name,
	}
}

// Backing returns the backing gocui view
func (view *View) Backing() *gocui.View {
	return view.backing
}

// SetBacking sets the backing gocui view
func (view *View) SetBacking(gocuiView *gocui.View) {
	view.backing = gocuiView
}

// HasBacking returns the true if this view has a gocui backing
func (view *View) HasBacking() bool {
	return view.backing != nil
}

// Size returns the view size
func (view *View) Size() (int, int) {
	if view.HasBacking() {
		return view.backing.Size()
	}

	return 0, 0
}

// Height returns the view height
func (view *View) Height() int {
	if view.HasBacking() {
		_, h := view.Size()
		return h
	}

	return 0
}

// Width returns the view width
func (view *View) Width() int {
	if view.HasBacking() {
		w, _ := view.Size()
		return w
	}

	return 0
}

// Cursor returns the view's cursor points
func (view *View) Cursor() (int, int) {
	if view.HasBacking() {
		return view.backing.Cursor()
	}

	return 0, 0
}

// CursorX returns the view's cursor X point
func (view *View) CursorX() int {
	if view.HasBacking() {
		x, _ := view.backing.Cursor()
		return x
	}

	return 0
}

// CursorY returns the view's cursor Y point
func (view *View) CursorY() int {
	if view.HasBacking() {
		_, y := view.backing.Cursor()
		return y
	}

	return 0
}

// SetCursor sets the view's cursor
func (view *View) SetCursor(x, y int) error {
	if view.HasBacking() {
		maxX, maxY := view.Size()
		if x < 0 || x >= maxX || y < 0 || y >= maxY {
			return nil
		}
		return view.backing.SetCursor(x, y)
	}

	return nil
}

// Origin returns the view's origin points
func (view *View) Origin() (int, int) {
	if view.HasBacking() {
		return view.backing.Origin()
	}

	return 0, 0
}

// OriginX returns the view's origin X point
func (view *View) OriginX() int {
	if view.HasBacking() {
		x, _ := view.backing.Origin()
		return x
	}

	return 0
}

// OriginY returns the view's origin Y point
func (view *View) OriginY() int {
	if view.HasBacking() {
		_, y := view.backing.Origin()
		return y
	}

	return 0
}

// SetOrigin sets the view's origin
func (view *View) SetOrigin(x, y int) error {
	if view.HasBacking() {
		if x < 0 || y < 0 {
			return nil
		}
		return view.backing.SetOrigin(x, y)
	}

	return nil
}

// Name returns the view's name
func (view *View) Name() string {
	return view.name
}

// Clear clears the view content
func (view *View) Clear() error {
	if view.HasBacking() {
		view.backing.Clear()
	}
	return nil
}

// Write will write the content to the view
func (view *View) Write(content string) error {
	if view.HasBacking() {
		fmt.Fprintln(view.backing, content)
	}
	return nil
}

// Update will clear and write the content to the view
func (view *View) Update(content string) error {
	view.Clear()
	view.Write(content)
	return nil
}

// SetFrame enables the frame border for the view
func (view *View) SetFrame(enabled bool) error {
	if view.HasBacking() {
		view.backing.Frame = enabled
	}
	return nil
}

// SetHighlight enables the highlight color for the view
func (view *View) SetHighlight(enabled bool) error {
	if view.HasBacking() {
		view.backing.Highlight = enabled
	}
	return nil
}

// SetEditable makes the view editable
func (view *View) SetEditable(enabled bool) error {
	if view.HasBacking() {
		view.backing.Editable = enabled
	}
	return nil
}

// SetWrap enables text wrapping for the view
func (view *View) SetWrap(enabled bool) error {
	if view.HasBacking() {
		view.backing.Wrap = enabled
	}
	return nil
}

// SetFgColor sets the foreground color
func (view *View) SetFgColor(color termbox.Attribute) {
	if view.HasBacking() {
		view.backing.Style = view.backing.Style.Foreground(termbox.MkColor(color))
	}
}

// SetBgColor sets the background color
func (view *View) SetBgColor(color termbox.Attribute) {
	if view.HasBacking() {
		// view.backing.BgColor = color
		view.backing.Style = view.backing.Style.Background(termbox.MkColor(color))
	}
}

// SetSelFgColor sets the foreground color for selection
func (view *View) SetSelFgColor(color termbox.Attribute) {
	if view.HasBacking() {
		view.backing.SelStyle = view.backing.SelStyle.Foreground(termbox.MkColor(color))
	}
}

// SetSelBgColor sets the background color for selection
func (view *View) SetSelBgColor(color termbox.Attribute) {
	if view.HasBacking() {
		view.backing.SelStyle = view.backing.SelStyle.Background(termbox.MkColor(color))
	}
}

// Read reads data in bytes buffer
func (view *View) Read(b []byte) (int, error) {
	if view.HasBacking() {
		return view.backing.Read(b)
	}
	return 0, nil
}

// Rewind undos view update
func (view *View) Rewind() error {
	if view.HasBacking() {
		view.backing.Rewind()
	}
	return nil
}
