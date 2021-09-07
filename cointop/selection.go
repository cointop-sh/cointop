package cointop

import log "github.com/sirupsen/logrus"

// SelectedCoinName returns the selected coin name
func (ct *Cointop) SelectedCoinName() string {
	log.Debug("SelectedCoinName()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Name
	}

	return ""
}

// SelectedCoinSymbol returns the selected coin symbol
func (ct *Cointop) SelectedCoinSymbol() string {
	log.Debug("SelectedCoinSymbol()")
	coin := ct.State.selectedCoin
	if coin != nil {
		return coin.Symbol
	}

	return ""
}
