package cointop

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/pad"
)

// PortfolioUpdateMenuView is structure for portfolio update menu view
type PortfolioUpdateMenuView struct {
	*View
}

// NewPortfolioUpdateMenuView returns a new portfolio update menu view
func NewPortfolioUpdateMenuView() *PortfolioUpdateMenuView {
	return &PortfolioUpdateMenuView{NewView("portfolioupdatemenu")}
}

// TogglePortfolio toggles the portfolio view
func (ct *Cointop) TogglePortfolio() error {
	ct.debuglog("togglePortfolio()")
	if ct.State.portfolioVisible {
		ct.GoToPageRowIndex(ct.State.lastSelectedRowIndex)
	} else {
		ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	}

	ct.State.filterByFavorites = false
	ct.State.portfolioVisible = !ct.State.portfolioVisible

	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

// ToggleShowPortfolio shows the portfolio view
func (ct *Cointop) ToggleShowPortfolio() error {
	ct.debuglog("toggleShowPortfolio()")
	ct.State.filterByFavorites = false
	ct.State.portfolioVisible = true
	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

// TogglePortfolioUpdateMenu toggles the portfolio update menu
func (ct *Cointop) TogglePortfolioUpdateMenu() error {
	ct.debuglog("togglePortfolioUpdateMenu()")
	ct.State.portfolioUpdateMenuVisible = !ct.State.portfolioUpdateMenuVisible
	if ct.State.portfolioUpdateMenuVisible {
		return ct.ShowPortfolioUpdateMenu()
	}

	return ct.HidePortfolioUpdateMenu()
}

// CoinHoldings returns portfolio coin holdings
func (ct *Cointop) CoinHoldings(coin *Coin) float64 {
	entry, _ := ct.PortfolioEntry(coin)
	return entry.Holdings
}

// UpdatePortfolioUpdateMenu updates the portfolio update menu view
func (ct *Cointop) UpdatePortfolioUpdateMenu() {
	ct.debuglog("updatePortfolioUpdateMenu()")
	coin := ct.HighlightedRowCoin()
	exists := ct.PortfolioEntryExists(coin)
	value := strconv.FormatFloat(ct.CoinHoldings(coin), 'f', -1, 64)
	ct.debuglog(fmt.Sprintf("holdings %v", value))
	var mode string
	var current string
	var submitText string
	if exists {
		mode = "Edit"
		current = fmt.Sprintf("(current %s %s)", value, coin.Symbol)
		submitText = "Set"
	} else {
		mode = "Add"
		submitText = "Add"
	}
	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" %s Portfolio Entry %s\n\n", mode, pad.Left("[q] close ", ct.maxTableWidth-26, " ")))
	label := fmt.Sprintf(" Enter holdings for %s %s", ct.colorscheme.MenuLabel(coin.Name), current)
	content := fmt.Sprintf("%s\n%s\n\n%s%s\n\n\n [Enter] %s    [ESC] Cancel", header, label, strings.Repeat(" ", 29), coin.Symbol, submitText)

	ct.Update(func() error {
		ct.Views.PortfolioUpdateMenu.Backing().Clear()
		ct.Views.PortfolioUpdateMenu.Backing().Frame = true
		fmt.Fprintln(ct.Views.PortfolioUpdateMenu.Backing(), content)
		fmt.Fprintln(ct.Views.Input.Backing(), value)
		ct.Views.Input.Backing().SetCursor(len(value), 0)
		return nil
	})
}

// ShowPortfolioUpdateMenu shows the portfolio update menu
func (ct *Cointop) ShowPortfolioUpdateMenu() error {
	ct.debuglog("showPortfolioUpdateMenu()")
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		ct.TogglePortfolio()
		return nil
	}

	ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	ct.State.portfolioUpdateMenuVisible = true
	ct.UpdatePortfolioUpdateMenu()
	ct.SetActiveView(ct.Views.PortfolioUpdateMenu.Name())
	return nil
}

// HidePortfolioUpdateMenu hides the portfolio update menu
func (ct *Cointop) HidePortfolioUpdateMenu() error {
	ct.debuglog("hidePortfolioUpdateMenu()")
	ct.State.portfolioUpdateMenuVisible = false
	ct.SetViewOnBottom(ct.Views.PortfolioUpdateMenu.Name())
	ct.SetViewOnBottom(ct.Views.Input.Name())
	ct.SetActiveView(ct.Views.Table.Name())
	ct.Update(func() error {
		if ct.Views.PortfolioUpdateMenu.Backing() == nil {
			return nil
		}

		ct.Views.PortfolioUpdateMenu.Backing().Clear()
		ct.Views.PortfolioUpdateMenu.Backing().Frame = false
		fmt.Fprintln(ct.Views.PortfolioUpdateMenu.Backing(), "")

		ct.Views.Input.Backing().Clear()
		fmt.Fprintln(ct.Views.Input.Backing(), "")
		return nil
	})

	return nil
}

// SetPortfolioHoldings sets portfolio entry holdings from inputed value
func (ct *Cointop) SetPortfolioHoldings() error {
	ct.debuglog("setPortfolioHoldings()")
	defer ct.HidePortfolioUpdateMenu()
	coin := ct.HighlightedRowCoin()

	// read input field
	b := make([]byte, 100)
	n, err := ct.Views.Input.Backing().Read(b)
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	value := normalizeFloatString(string(b))
	shouldDelete := value == ""
	var holdings float64

	if !shouldDelete {
		holdings, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	}

	if err := ct.SetPortfolioEntry(coin.Name, holdings); err != nil {
		return err
	}

	if shouldDelete {
		ct.RemovePortfolioEntry(coin.Name)
		ct.UpdateTable()
	} else {
		ct.UpdateTable()
		ct.GoToPageRowIndex(ct.State.lastSelectedRowIndex)
	}

	return nil
}

// PortfolioEntry returns a portfolio entry
func (ct *Cointop) PortfolioEntry(c *Coin) (*PortfolioEntry, bool) {
	//ct.debuglog("portfolioEntry()")
	if c == nil {
		return &PortfolioEntry{}, true
	}

	var p *PortfolioEntry
	var isNew bool
	var ok bool
	key := strings.ToLower(c.Name)
	if p, ok = ct.State.portfolio.Entries[key]; !ok {
		// NOTE: if not found then try the symbol
		key := strings.ToLower(c.Symbol)
		if p, ok = ct.State.portfolio.Entries[key]; !ok {
			p = &PortfolioEntry{
				Coin:     c.Name,
				Holdings: 0,
			}
			isNew = true
		}
	}

	return p, isNew
}

// SetPortfolioEntry sets a portfolio entry
func (ct *Cointop) SetPortfolioEntry(coin string, holdings float64) error {
	ct.debuglog("setPortfolioEntry()")
	ic, _ := ct.State.allCoinsSlugMap.Load(strings.ToLower(coin))
	c, _ := ic.(*Coin)
	p, isNew := ct.PortfolioEntry(c)
	if isNew {
		key := strings.ToLower(coin)
		ct.State.portfolio.Entries[key] = &PortfolioEntry{
			Coin:     coin,
			Holdings: holdings,
		}
	} else {
		p.Holdings = holdings
	}

	if err := ct.Save(); err != nil {
		return err
	}

	return nil
}

// RemovePortfolioEntry removes a portfolio entry
func (ct *Cointop) RemovePortfolioEntry(coin string) {
	ct.debuglog("removePortfolioEntry()")
	delete(ct.State.portfolio.Entries, strings.ToLower(coin))
}

// PortfolioEntryExists returns true if portfolio entry exists
func (ct *Cointop) PortfolioEntryExists(c *Coin) bool {
	ct.debuglog("portfolioEntryExists()")
	_, isNew := ct.PortfolioEntry(c)
	return !isNew
}

// PortfolioEntriesCount returns the count of portfolio entries
func (ct *Cointop) PortfolioEntriesCount() int {
	ct.debuglog("portfolioEntriesCount()")
	return len(ct.State.portfolio.Entries)
}

// GetPortfolioSlice returns portfolio entries as a slice
func (ct *Cointop) GetPortfolioSlice() []*Coin {
	ct.debuglog("getPortfolioSlice()")
	sliced := []*Coin{}
	if ct.PortfolioEntriesCount() == 0 {
		return sliced
	}

	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		p, isNew := ct.PortfolioEntry(coin)
		if isNew {
			continue
		}
		coin.Holdings = p.Holdings
		balance := coin.Price * p.Holdings
		balancestr := fmt.Sprintf("%.2f", balance)
		if ct.State.currencyConversion == "ETH" || ct.State.currencyConversion == "BTC" {
			balancestr = fmt.Sprintf("%.5f", balance)
		}
		balance, _ = strconv.ParseFloat(balancestr, 64)
		coin.Balance = balance
		sliced = append(sliced, coin)
	}

	sort.Slice(sliced, func(i, j int) bool {
		return sliced[i].Balance > sliced[j].Balance
	})

	for i, coin := range sliced {
		coin.Rank = i + 1
	}

	return sliced
}

// GetPortfolioTotal returns the total balance of portfolio entries
func (ct *Cointop) GetPortfolioTotal() float64 {
	ct.debuglog("getPortfolioTotal()")
	portfolio := ct.GetPortfolioSlice()
	var total float64
	for _, p := range portfolio {
		total += p.Balance
	}
	return total
}

// NormalizeFloatString normalizes a float as a string
func normalizeFloatString(input string) string {
	re := regexp.MustCompile(`(\d+\.\d+|\.\d+|\d+)`)
	result := re.FindStringSubmatch(input)
	if len(result) > 0 {
		return result[0]
	}

	return ""
}
