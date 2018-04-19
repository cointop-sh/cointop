package cointop

// Size returns window width and height
func (ct *Cointop) size() (int, int) {
	return ct.g.Size()
}

// Width returns window width
func (ct *Cointop) width() int {
	w, _ := ct.size()
	return w
}

// Height returns window height
func (ct *Cointop) height() int {
	_, h := ct.size()
	return h
}
