package cointop

const version = "1.0.6"

func (ct *Cointop) version() string {
	return version
}

// Version returns cointop version
func Version() string {
	return version
}
