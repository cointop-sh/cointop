package cointop

const version = "1.1.1"

func (ct *Cointop) version() string {
	return version
}

// Version returns cointop version
func Version() string {
	return version
}
