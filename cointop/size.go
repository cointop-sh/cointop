package cointop

// Size returns window width and height
func (ct *Cointop) Size() (int, int) {
	return ct.g.Size()
}

// Width returns window width
func (ct *Cointop) Width() int {
	w, _ := ct.Size()
	return w
}

// Height returns window height
func (ct *Cointop) Height() int {
	_, h := ct.Size()
	return h
}
