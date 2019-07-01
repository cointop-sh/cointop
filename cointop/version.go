package cointop

// TODO: make dynamic based on git tag
const version = "1.3.2"

// Version returns the cointop version
func (ct *Cointop) Version() string {
	return version
}

// Version returns cointop version
func Version() string {
	return version
}
