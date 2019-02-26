package cointop

// TODO: make dynamic based on git tag
const version = "1.1.3"

func (ct *Cointop) version() string {
	return version
}

// Version returns cointop version
func Version() string {
	return version
}
