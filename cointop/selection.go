package cointop

func (ct *Cointop) selectedCoinName() string {
	ct.debuglog("selectedCoinName()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

func (ct *Cointop) selectedCoinSymbol() string {
	ct.debuglog("selectedCoinSymbol()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Symbol
	}

	return ""
}
