package cointop

// Size returns window width and height
func (ct *Cointop) size() (int, int) {
	ct.debuglog("size()")
	if ct.g == nil {
		return 0, 0
	}

	return ct.g.Size()
}

// Width returns window width
func (ct *Cointop) width() int {
	ct.debuglog("width()")
	w, _ := ct.size()
	return w
}

// Height returns window height
func (ct *Cointop) height() int {
	ct.debuglog("height()")
	_, h := ct.size()
	return h
}

// ViewWidth returns view width
func (ct *Cointop) ViewWidth(view string) int {
	ct.debuglog("viewWidth()")
	v, err := ct.g.View(view)
	if err != nil {
		return 0
	}
	w, _ := v.Size()
	return w
}

// ClampedWidth returns the clamped width
func (ct *Cointop) ClampedWidth() int {
	ct.debuglog("clampedWidth()")
	w := ct.width()
	if w > ct.maxTableWidth {
		return ct.maxTableWidth
	}

	return w
}
