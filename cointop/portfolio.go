package cointop

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/pad"
	log "github.com/sirupsen/logrus"
)

// PortfolioUpdateMenuView is structure for portfolio update menu view
type PortfolioUpdateMenuView struct {
	*View
}

// NewPortfolioUpdateMenuView returns a new portfolio update menu view
func NewPortfolioUpdateMenuView() *PortfolioUpdateMenuView {
	return &PortfolioUpdateMenuView{NewView("portfolioupdatemenu")}
}

func (ct *Cointop) togglePortfolio() error {
	ct.debuglog("togglePortfolio()")
	if ct.State.portfolioVisible {
		ct.goToPageRowIndex(ct.State.lastSelectedRowIndex)
	} else {
		ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	}

	ct.State.filterByFavorites = false
	ct.State.portfolioVisible = !ct.State.portfolioVisible

	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

func (ct *Cointop) toggleShowPortfolio() error {
	ct.debuglog("toggleShowPortfolio()")
	ct.State.filterByFavorites = false
	ct.State.portfolioVisible = true
	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

func (ct *Cointop) togglePortfolioUpdateMenu() error {
	ct.debuglog("togglePortfolioUpdateMenu()")
	ct.State.portfolioUpdateMenuVisible = !ct.State.portfolioUpdateMenuVisible
	if ct.State.portfolioUpdateMenuVisible {
		return ct.showPortfolioUpdateMenu()
	}

	return ct.hidePortfolioUpdateMenu()
}

func (ct *Cointop) coinHoldings(coin *Coin) float64 {
	entry, _ := ct.PortfolioEntry(coin)
	return entry.Holdings
}

func (ct *Cointop) updatePortfolioUpdateMenu() {
	ct.debuglog("updatePortfolioUpdateMenu()")
	coin := ct.HighlightedRowCoin()
	exists := ct.PortfolioEntryExists(coin)
	value := strconv.FormatFloat(ct.coinHoldings(coin), 'f', -1, 64)
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

	ct.Update(func() {
		ct.Views.PortfolioUpdateMenu.Backing().Clear()
		ct.Views.PortfolioUpdateMenu.Backing().Frame = true
		fmt.Fprintln(ct.Views.PortfolioUpdateMenu.Backing(), content)
		fmt.Fprintln(ct.Views.Input.Backing(), value)
		ct.Views.Input.Backing().SetCursor(len(value), 0)
	})
}

func (ct *Cointop) showPortfolioUpdateMenu() error {
	ct.debuglog("showPortfolioUpdateMenu()")
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		ct.togglePortfolio()
		return nil
	}

	ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	ct.State.portfolioUpdateMenuVisible = true
	ct.updatePortfolioUpdateMenu()
	ct.SetActiveView(ct.Views.PortfolioUpdateMenu.Name())
	return nil
}

func (ct *Cointop) hidePortfolioUpdateMenu() error {
	ct.debuglog("hidePortfolioUpdateMenu()")
	ct.State.portfolioUpdateMenuVisible = false
	ct.SetViewOnBottom(ct.Views.PortfolioUpdateMenu.Name())
	ct.SetViewOnBottom(ct.Views.Input.Name())
	ct.SetActiveView(ct.Views.Table.Name())
	ct.Update(func() {
		if ct.Views.PortfolioUpdateMenu.Backing() == nil {
			return
		}

		ct.Views.PortfolioUpdateMenu.Backing().Clear()
		ct.Views.PortfolioUpdateMenu.Backing().Frame = false
		fmt.Fprintln(ct.Views.PortfolioUpdateMenu.Backing(), "")

		ct.Views.Input.Backing().Clear()
		fmt.Fprintln(ct.Views.Input.Backing(), "")
	})
	return nil
}

// sets portfolio entry holdings from inputed value
func (ct *Cointop) setPortfolioHoldings() error {
	ct.debuglog("setPortfolioHoldings()")
	defer ct.hidePortfolioUpdateMenu()
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

	value := normalizeFloatstring(string(b))
	shouldDelete := value == ""
	var holdings float64

	if !shouldDelete {
		holdings, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	}

	ct.setPortfolioEntry(coin.Name, holdings)

	if shouldDelete {
		ct.removePortfolioEntry(coin.Name)
		ct.UpdateTable()
	} else {
		ct.UpdateTable()
		ct.goToPageRowIndex(ct.State.lastSelectedRowIndex)
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

func (ct *Cointop) setPortfolioEntry(coin string, holdings float64) {
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
		log.Fatal(err)
	}
}

func (ct *Cointop) removePortfolioEntry(coin string) {
	ct.debuglog("removePortfolioEntry()")
	delete(ct.State.portfolio.Entries, strings.ToLower(coin))
}

// PortfolioEntryExists returns true if portfolio entry exists
func (ct *Cointop) PortfolioEntryExists(c *Coin) bool {
	ct.debuglog("portfolioEntryExists()")
	_, isNew := ct.PortfolioEntry(c)
	return !isNew
}

func (ct *Cointop) portfolioEntriesCount() int {
	ct.debuglog("portfolioEntriesCount()")
	return len(ct.State.portfolio.Entries)
}

func (ct *Cointop) getPortfolioSlice() []*Coin {
	ct.debuglog("getPortfolioSlice()")
	sliced := []*Coin{}
	if ct.portfolioEntriesCount() == 0 {
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

func (ct *Cointop) getPortfolioTotal() float64 {
	ct.debuglog("getPortfolioTotal()")
	portfolio := ct.getPortfolioSlice()
	var total float64
	for _, p := range portfolio {
		total += p.Balance
	}
	return total
}

func normalizeFloatstring(input string) string {
	re := regexp.MustCompile(`(\d+\.\d+|\.\d+|\d+)`)
	result := re.FindStringSubmatch(input)
	if len(result) > 0 {
		return result[0]
	}

	return ""
}
