package cointop

func (ct *Cointop) selectedCoinName() string {
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

func (ct *Cointop) selectedCoinSymbol() string {
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Symbol
	}

	return ""
}
