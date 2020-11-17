package cointop

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/miguelmota/cointop/pkg/asciitable"
	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/ui"
)

// PortfolioUpdateMenuView is structure for portfolio update menu view
type PortfolioUpdateMenuView = ui.View

// NewPortfolioUpdateMenuView returns a new portfolio update menu view
func NewPortfolioUpdateMenuView() *PortfolioUpdateMenuView {
	var view *PortfolioUpdateMenuView = ui.NewView("portfolioupdatemenu")
	return view
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
func (ct *Cointop) UpdatePortfolioUpdateMenu() error {
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

	ct.UpdateUI(func() error {
		ct.Views.PortfolioUpdateMenu.SetFrame(true)
		ct.Views.PortfolioUpdateMenu.Update(content)
		ct.Views.Input.Write(value)
		ct.Views.Input.SetCursor(len(value), 0)
		return nil
	})
	return nil
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
	ct.ui.SetViewOnBottom(ct.Views.PortfolioUpdateMenu)
	ct.ui.SetViewOnBottom(ct.Views.Input)
	ct.SetActiveView(ct.Views.Table.Name())
	ct.UpdateUI(func() error {
		ct.Views.PortfolioUpdateMenu.SetFrame(false)
		ct.Views.PortfolioUpdateMenu.Update("")
		ct.Views.Input.Update("")
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
	n, err := ct.Views.Input.Read(b)
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
	//ct.debuglog("portfolioEntry()") // too many
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

// RefreshPortfolioCoins refreshes portfolio entry coin data
func (ct *Cointop) RefreshPortfolioCoins() error {
	ct.debuglog("refreshPortfolioCoins()")
	holdings := ct.GetPortfolioSlice()
	holdingCoins := make([]string, len(holdings))
	for i, entry := range holdings {
		holdingCoins[i] = entry.Name
	}

	coins, err := ct.api.GetCoinDataBatch(holdingCoins, ct.State.currencyConversion)
	ct.processCoins(coins)
	if err != nil {
		return err
	}

	return nil
}

// TablePrintOptions are options for ascii table output.
type TablePrintOptions struct {
	SortBy        string
	SortDesc      bool
	HumanReadable bool
	Format        string
	Filter        []string
	Convert       string
}

// outputFormats is list of valid output formats
var outputFormats = map[string]bool{
	"table": true,
	"csv":   true,
	"json":  true,
}

// portfolioColumns is list of valid column keys for portfolio
var portfolioColumns = map[string]bool{
	"name":     true,
	"symbol":   true,
	"price":    true,
	"holdings": true,
	"balance":  true,
	"24h":      true,
}

// PrintHoldingsTable prints the holdings in an ASCII table
func (ct *Cointop) PrintHoldingsTable(options *TablePrintOptions) error {
	ct.debuglog("printHoldingsTable()")
	if options == nil {
		options = &TablePrintOptions{}
	}

	if err := ct.SetCurrencyConverstion(options.Convert); err != nil {
		return err
	}

	ct.RefreshPortfolioCoins()

	sortBy := options.SortBy
	sortDesc := options.SortDesc
	format := options.Format
	humanReadable := options.HumanReadable
	filter := options.Filter
	holdings := ct.GetPortfolioSlice()

	if format == "" {
		format = "table"
	}

	if sortBy != "" {
		if _, ok := portfolioColumns[sortBy]; !ok {
			return fmt.Errorf("The option %q is not a valid column name", sortBy)
		}

		ct.Sort(sortBy, sortDesc, holdings, true)
	}

	if _, ok := outputFormats[format]; !ok {
		return fmt.Errorf("The option %q is not a valid format type", format)
	}

	total := ct.GetPortfolioTotal()
	records := make([][]string, len(holdings))
	symbol := ct.CurrencySymbol()

	for i, entry := range holdings {
		if filter != nil && len(filter) > 0 {
			found := false
			for _, item := range filter {
				item = strings.ToLower(strings.TrimSpace(item))
				if strings.ToLower(entry.Symbol) == item || strings.ToLower(entry.Name) == item {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		percentHoldings := (entry.Balance / total) * 1e2
		if math.IsNaN(percentHoldings) {
			percentHoldings = 0
		}

		if humanReadable {
			records[i] = []string{
				entry.Name,
				entry.Symbol,
				fmt.Sprintf("%s%s", symbol, humanize.Commaf(entry.Price)),
				humanize.Commaf(entry.Holdings),
				fmt.Sprintf("%s%s", symbol, humanize.Commaf(entry.Balance)),
				fmt.Sprintf("%.2f%%", entry.PercentChange24H),
				fmt.Sprintf("%.2f%%", percentHoldings),
			}
		} else {
			records[i] = []string{
				entry.Name,
				entry.Symbol,
				strconv.FormatFloat(entry.Price, 'f', -1, 64),
				strconv.FormatFloat(entry.Holdings, 'f', -1, 64),
				strconv.FormatFloat(entry.Balance, 'f', -1, 64),
				fmt.Sprintf("%.2f", entry.PercentChange24H),
				fmt.Sprintf("%.2f", percentHoldings),
			}
		}
	}

	headers := []string{"name", "symbol", "price", "holdings", "balance", "24h%", "%holdings"}

	if format == "csv" {
		csvWriter := csv.NewWriter(os.Stdout)
		if err := csvWriter.Write(headers); err != nil {
			return err
		}

		for _, record := range records {
			if err := csvWriter.Write(record); err != nil {
				return err
			}
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			return err
		}

		return nil
	} else if format == "json" {
		list := make([]map[string]string, len(records))
		for i, record := range records {
			obj := make(map[string]string, len(record))
			for j, column := range record {
				obj[headers[j]] = column
			}

			list[i] = obj
		}

		output, err := json.Marshal(list)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	}

	alignment := []int{-1, -1, 1, 1, 1, 1, 1}
	table := asciitable.NewAsciiTable(&asciitable.Input{
		Data:      records,
		Headers:   headers,
		Alignment: alignment,
	})

	fmt.Println(table.String())
	return nil
}

// PrintTotalHoldings prints the total holdings amount
func (ct *Cointop) PrintTotalHoldings(options *TablePrintOptions) error {
	ct.debuglog("printTotalHoldings()")
	if options == nil {
		options = &TablePrintOptions{}
	}

	if err := ct.SetCurrencyConverstion(options.Convert); err != nil {
		return err
	}

	ct.RefreshPortfolioCoins()

	humanReadable := options.HumanReadable
	symbol := ct.CurrencySymbol()
	format := options.Format
	filter := options.Filter
	portfolio := ct.GetPortfolioSlice()
	var total float64
	for _, entry := range portfolio {
		if filter != nil && len(filter) > 0 {
			found := false
			for _, item := range filter {
				item = strings.ToLower(strings.TrimSpace(item))
				if strings.ToLower(entry.Symbol) == item || strings.ToLower(entry.Name) == item {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		total += entry.Balance
	}

	value := strconv.FormatFloat(total, 'f', -1, 64)

	if humanReadable {
		value = fmt.Sprintf("%s%s", symbol, humanize.Commaf(total))
	}

	if format == "csv" {
		csvWriter := csv.NewWriter(os.Stdout)
		if err := csvWriter.Write([]string{"total"}); err != nil {
			return err
		}
		if err := csvWriter.Write([]string{value}); err != nil {
			return err
		}

		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			return err
		}

		return nil
	} else if format == "json" {
		obj := map[string]string{"total": value}
		output, err := json.Marshal(obj)
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	}

	fmt.Println(value)

	return nil
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
