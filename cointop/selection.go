package cointop

// SelectedCoinName returns the selected coin name
func (ct *Cointop) SelectedCoinName() string {
	ct.debuglog("selectedCoinName()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

// SelectedCoinSymbol returns the selected coin symbol
func (ct *Cointop) SelectedCoinSymbol() string {
	ct.debuglog("selectedCoinSymbol()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Symbol
	}

	return ""
}
