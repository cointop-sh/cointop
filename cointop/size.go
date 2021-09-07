package cointop

import log "github.com/sirupsen/logrus"

// Size returns window width and height
func (ct *Cointop) Size() (int, int) {
	log.Debug("Size()")
	if ct.g == nil {
		return 0, 0
	}

	return ct.g.Size()
}

// Width returns window width
func (ct *Cointop) Width() int {
	log.Debug("Width()")
	w, _ := ct.Size()
	return w
}

// Height returns window height
func (ct *Cointop) Height() int {
	log.Debug("Height()")
	_, h := ct.Size()
	return h
}

// ViewWidth returns view width
func (ct *Cointop) ViewWidth(view string) int {
	log.Debug("ViewWidth()")
	v, err := ct.g.View(view)
	if err != nil {
		return 0
	}
	w, _ := v.Size()
	return w
}

// ClampedWidth returns the clamped width
func (ct *Cointop) ClampedWidth() int {
	log.Debug("ClampedWidth()")
	w := ct.Width()
	if w > ct.maxTableWidth {
		return ct.maxTableWidth
	}

	return w
}
