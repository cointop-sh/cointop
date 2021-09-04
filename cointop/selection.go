package cointop

// SelectedCoinName returns the selected coin name
func (ct *Cointop) SelectedCoinName() string {
	ct.debuglog("SelectedCoinName()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

// SelectedCoinSymbol returns the selected coin symbol
func (ct *Cointop) SelectedCoinSymbol() string {
	ct.debuglog("SelectedCoinSymbol()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Symbol
	}

	return ""
}
