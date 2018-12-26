package cointop

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/pad"
)

func (ct *Cointop) togglePortfolio() error {
	ct.filterByFavorites = false
	ct.portfoliovisible = !ct.portfoliovisible
	ct.updateTable()
	return nil
}

func (ct *Cointop) toggleShowPortfolio() error {
	ct.filterByFavorites = false
	ct.portfoliovisible = true
	ct.updateTable()
	return nil
}

func (ct *Cointop) togglePortfolioUpdateMenu() error {
	ct.portfolioupdatemenuvisible = !ct.portfolioupdatemenuvisible
	if ct.portfolioupdatemenuvisible {
		return ct.showPortfolioUpdateMenu()
	}
	return ct.hidePortfolioUpdateMenu()
}

func (ct *Cointop) updatePortfolioUpdateMenu() {
	coin := ct.highlightedRowCoin()
	exists := ct.portfolioEntryExists(coin)
	value := strconv.FormatFloat(coin.Holdings, 'f', -1, 64)
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
	header := color.GreenBg(fmt.Sprintf(" %s Portfolio Entry %s\n\n", mode, pad.Left("[q] close ", ct.maxtablewidth-26, " ")))
	label := fmt.Sprintf(" Enter holdings for %s %s", color.Yellow(coin.Name), current)
	content := fmt.Sprintf("%s\n%s\n\n%s%s\n\n\n [Enter] %s    [ESC] Cancel", header, label, strings.Repeat(" ", 29), coin.Symbol, submitText)

	ct.update(func() {
		ct.portfolioupdatemenuview.Clear()
		ct.portfolioupdatemenuview.Frame = true
		fmt.Fprintln(ct.portfolioupdatemenuview, content)
		fmt.Fprintln(ct.inputview, value)
		ct.inputview.SetCursor(len(value), 0)
	})
}

func (ct *Cointop) showPortfolioUpdateMenu() error {
	ct.portfolioupdatemenuvisible = true
	ct.updatePortfolioUpdateMenu()
	ct.setActiveView(ct.portfolioupdatemenuviewname)
	return nil
}

func (ct *Cointop) hidePortfolioUpdateMenu() error {
	ct.portfolioupdatemenuvisible = false
	ct.setViewOnBottom(ct.portfolioupdatemenuviewname)
	ct.setViewOnBottom(ct.inputviewname)
	ct.setActiveView(ct.tableviewname)
	ct.update(func() {
		ct.portfolioupdatemenuview.Clear()
		ct.portfolioupdatemenuview.Frame = false
		fmt.Fprintln(ct.portfolioupdatemenuview, "")

		ct.inputview.Clear()
		fmt.Fprintln(ct.inputview, "")
	})
	return nil
}

// sets portfolio entry holdings from inputed value
func (ct *Cointop) setPortfolioHoldings() error {
	defer ct.hidePortfolioUpdateMenu()
	coin := ct.highlightedRowCoin()

	b := make([]byte, 100)
	n, err := ct.inputview.Read(b)
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
		ct.updateTable()
		ct.goToGlobalIndex(0)
	} else {
		ct.updateTable()
		ct.goToGlobalIndex(coin.Rank - 1)
	}

	return nil
}

func (ct *Cointop) portfolioEntry(c *coin) (*portfolioEntry, bool) {
	if c == nil {
		return &portfolioEntry{}, true
	}

	var p *portfolioEntry
	var isNew bool
	var ok bool
	key := strings.ToLower(c.Name)
	if p, ok = ct.portfolio.Entries[key]; !ok {
		// NOTE: if not found then try the symbol
		key := strings.ToLower(c.Symbol)
		if p, ok = ct.portfolio.Entries[key]; !ok {
			p = &portfolioEntry{
				Coin:     c.Name,
				Holdings: 0,
			}
			isNew = true
		}
	}

	return p, isNew
}

func (ct *Cointop) setPortfolioEntry(coin string, holdings float64) {
	c, _ := ct.allcoinsslugmap[strings.ToLower(coin)]
	p, isNew := ct.portfolioEntry(c)
	if isNew {
		key := strings.ToLower(coin)
		ct.portfolio.Entries[key] = &portfolioEntry{
			Coin:     coin,
			Holdings: holdings,
		}
	} else {
		p.Holdings = holdings
	}
}

func (ct *Cointop) removePortfolioEntry(coin string) {
	delete(ct.portfolio.Entries, strings.ToLower(coin))
}

func (ct *Cointop) portfolioEntryExists(c *coin) bool {
	_, isNew := ct.portfolioEntry(c)
	return !isNew
}

func (ct *Cointop) portfolioEntriesCount() int {
	return len(ct.portfolio.Entries)
}

func (ct *Cointop) getPortfolioSlice() []*coin {
	sliced := []*coin{}
	for i := range ct.allcoins {
		if ct.portfolioEntriesCount() == 0 {
			break
		}
		coin := ct.allcoins[i]
		p, isNew := ct.portfolioEntry(coin)
		if isNew {
			continue
		}
		coin.Holdings = p.Holdings
		balance := coin.Price * p.Holdings
		balancestr := fmt.Sprintf("%.2f", balance)
		if ct.currencyconversion == "ETH" || ct.currencyconversion == "BTC" {
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
