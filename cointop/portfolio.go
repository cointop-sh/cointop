package cointop

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cointop-sh/cointop/pkg/asciitable"
	"github.com/cointop-sh/cointop/pkg/eval"
	"github.com/cointop-sh/cointop/pkg/humanize"
	"github.com/cointop-sh/cointop/pkg/pad"
	"github.com/cointop-sh/cointop/pkg/table"
	log "github.com/sirupsen/logrus"
)

// SupportedPortfolioTableHeaders are all the supported portfolio table header columns
var SupportedPortfolioTableHeaders = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"holdings",
	"balance",
	"1h_change",
	"24h_change",
	"7d_change",
	"30d_change",
	"1y_change",
	"percent_holdings",
	"last_updated",
}

// DefaultPortfolioTableHeaders are the default portfolio table header columns
var DefaultPortfolioTableHeaders = []string{
	"rank",
	"name",
	"symbol",
	"price",
	"holdings",
	"balance",
	"1h_change",
	"24h_change",
	"7d_change",
	"percent_holdings",
	"last_updated",
}

// HiddenBalanceChars are the characters to show when hidding balances
var HiddenBalanceChars = "********"

// ValidPortfolioTableHeader returns the portfolio table headers
func (ct *Cointop) ValidPortfolioTableHeader(name string) bool {
	for _, v := range SupportedPortfolioTableHeaders {
		if v == name {
			return true
		}
	}

	return false
}

// GetPortfolioTableHeaders returns the portfolio table headers
func (ct *Cointop) GetPortfolioTableHeaders() []string {
	return ct.State.portfolioTableColumns
}

// GetPortfolioTable returns the table for displaying portfolio holdings
func (ct *Cointop) GetPortfolioTable() *table.Table {
	total := ct.GetPortfolioTotal()
	maxX := ct.Width()
	t := table.NewTable().SetWidth(maxX)
	var rows [][]*table.RowCell
	headers := ct.GetPortfolioTableHeaders()
	ct.ClearSyncMap(&ct.State.tableColumnWidths)
	ct.ClearSyncMap(&ct.State.tableColumnAlignLeft)
	for _, coin := range ct.State.coins {
		leftMargin := 1
		rightMargin := 1
		var rowCells []*table.RowCell
		for _, header := range headers {
			switch header {
			case "rank":
				star := ct.colorscheme.TableRow(" ")
				if coin.Favorite {
					star = ct.colorscheme.TableRowFavorite(ct.State.favoriteChar)
				}
				rank := fmt.Sprintf("%s%v", star, ct.colorscheme.TableRow(fmt.Sprintf("%6v ", coin.Rank)))
				ct.SetTableColumnWidth(header, 8)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       ct.colorscheme.Default,
					Text:        rank,
				})
			case "name":

				name := TruncateString(coin.Name, 18)
				namecolor := ct.colorscheme.TableRow
				if coin.Favorite {
					namecolor = ct.colorscheme.TableRowFavorite
				}
				ct.SetTableColumnWidthFromString(header, name)
				ct.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   true,
						Color:       namecolor,
						Text:        name,
					})
			case "symbol":
				symbol := TruncateString(coin.Symbol, 6)
				ct.SetTableColumnWidthFromString(header, symbol)
				ct.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   true,
						Color:       ct.colorscheme.TableRow,
						Text:        symbol,
					})
			case "price":
				text := ct.FormatPrice(coin.Price)
				symbolPadding := 1
				ct.SetTableColumnWidth(header, utf8.RuneCountInString(text)+symbolPadding)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "holdings":
				text := strconv.FormatFloat(coin.Holdings, 'f', -1, 64)
				if ct.State.hidePortfolioBalances {
					text = HiddenBalanceChars
				}
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "balance":
				text := humanize.Monetaryf(coin.Balance, 2)
				if ct.State.hidePortfolioBalances {
					text = HiddenBalanceChars
				}
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				colorBalance := ct.colorscheme.TableColumnPrice
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       colorBalance,
						Text:        text,
					})
			case "1h_change":
				color1h := ct.colorscheme.TableColumnChange
				if coin.PercentChange1H > 0 {
					color1h = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange1H < 0 {
					color1h = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange1H)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color1h,
						Text:        text,
					})
			case "24h_change":
				color24h := ct.colorscheme.TableColumnChange
				if coin.PercentChange24H > 0 {
					color24h = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange24H < 0 {
					color24h = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange24H)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color24h,
						Text:        text,
					})
			case "7d_change":
				color7d := ct.colorscheme.TableColumnChange
				if coin.PercentChange7D > 0 {
					color7d = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange7D < 0 {
					color7d = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange7D)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color7d,
						Text:        text,
					})
			case "30d_change":
				color30d := ct.colorscheme.TableColumnChange
				if coin.PercentChange30D > 0 {
					color30d = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange30D < 0 {
					color30d = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange30D)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color30d,
						Text:        text,
					})
			case "1y_change":
				color1y := ct.colorscheme.TableColumnChange
				if coin.PercentChange1Y > 0 {
					color1y = ct.colorscheme.TableColumnChangeUp
				}
				if coin.PercentChange1Y < 0 {
					color1y = ct.colorscheme.TableColumnChangeDown
				}
				text := fmt.Sprintf("%.2f%%", coin.PercentChange1Y)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       color1y,
						Text:        text,
					})
			case "percent_holdings":
				percentHoldings := (coin.Balance / total) * 1e2
				if math.IsNaN(percentHoldings) {
					percentHoldings = 0
				}
				text := fmt.Sprintf("%.2f%%", percentHoldings)
				if ct.State.hidePortfolioBalances {
					text = HiddenBalanceChars
				}
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        text,
					})
			case "last_updated":
				unix, _ := strconv.ParseInt(coin.LastUpdated, 10, 64)
				lastUpdated := time.Unix(unix, 0).Format("15:04:05 Jan 02")
				ct.SetTableColumnWidthFromString(header, lastUpdated)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells,
					&table.RowCell{
						LeftMargin:  leftMargin,
						RightMargin: rightMargin,
						LeftAlign:   false,
						Color:       ct.colorscheme.TableRow,
						Text:        lastUpdated,
					})
			}
		}

		rows = append(rows, rowCells)
	}

	for _, row := range rows {
		for i, header := range headers {
			row[i].Width = ct.GetTableColumnWidth(header)
		}
		t.AddRowCells(row...)
	}

	return t
}

// TogglePortfolio toggles the portfolio view
func (ct *Cointop) TogglePortfolio() error {
	log.Debug("TogglePortfolio()")
	ct.ToggleSelectedView(PortfolioView)
	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

// ToggleShowPortfolio shows the portfolio view
func (ct *Cointop) ToggleShowPortfolio() error {
	log.Debug("ToggleShowPortfolio()")
	ct.SetSelectedView(PortfolioView)
	go ct.UpdateChart()
	go ct.UpdateTable()
	return nil
}

// TogglePortfolioUpdateMenu toggles the portfolio update menu
func (ct *Cointop) TogglePortfolioUpdateMenu() error {
	log.Debug("TogglePortfolioUpdateMenu()")
	if ct.IsPriceAlertsVisible() {
		return ct.ShowPriceAlertsUpdateMenu()
	}
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
	log.Debug("UpdatePortfolioUpdateMenu()")
	coin := ct.HighlightedRowCoin()
	exists := ct.PortfolioEntryExists(coin)
	value := strconv.FormatFloat(ct.CoinHoldings(coin), 'f', -1, 64)
	log.Debugf("UpdatePortfolioUpdateMenu() holdings %v", value)
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
	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" %s Portfolio Entry %s\n\n", mode, pad.Left("[q] close ", ct.Width()-25, " ")))
	label := fmt.Sprintf(" Enter holdings for %s %s", ct.colorscheme.MenuLabel(coin.Name), current)
	content := fmt.Sprintf("%s\n%s\n\n%s%s\n\n\n [Enter] %s    [ESC] Cancel", header, label, strings.Repeat(" ", 29), coin.Symbol, submitText)

	ct.UpdateUI(func() error {
		ct.Views.Menu.SetFrame(true)
		ct.Views.Menu.Update(content)
		ct.Views.Input.Write(value)
		ct.Views.Input.SetCursor(utf8.RuneCountInString(value), 0)
		return nil
	})
	return nil
}

// ShowPortfolioUpdateMenu shows the portfolio update menu
func (ct *Cointop) ShowPortfolioUpdateMenu() error {
	log.Debug("ShowPortfolioUpdateMenu()")

	// TODO: separation of concerns
	if ct.IsPriceAlertsVisible() {
		return ct.ShowPriceAlertsUpdateMenu()
	}

	coin := ct.HighlightedRowCoin()
	if coin == nil {
		ct.TogglePortfolio()
		return nil
	}

	ct.State.portfolioUpdateMenuVisible = true
	ct.UpdatePortfolioUpdateMenu()
	ct.ui.SetCursor(true)
	ct.SetActiveView(ct.Views.Menu.Name())
	ct.g.SetViewOnTop(ct.Views.Input.Name())
	ct.g.SetCurrentView(ct.Views.Input.Name())
	return nil
}

// HidePortfolioUpdateMenu hides the portfolio update menu
func (ct *Cointop) HidePortfolioUpdateMenu() error {
	log.Debug("HidePortfolioUpdateMenu()")
	ct.State.portfolioUpdateMenuVisible = false
	ct.ui.SetViewOnBottom(ct.Views.Menu)
	ct.ui.SetViewOnBottom(ct.Views.Input)
	ct.ui.SetCursor(false)
	ct.SetActiveView(ct.Views.Table.Name())
	ct.UpdateUI(func() error {
		ct.Views.Menu.SetFrame(false)
		ct.Views.Menu.Update("")
		ct.Views.Input.Update("")
		return nil
	})

	return nil
}

// SetPortfolioHoldings sets portfolio entry holdings from inputed value
func (ct *Cointop) SetPortfolioHoldings() error {
	log.Debug("SetPortfolioHoldings()")
	defer ct.HidePortfolioUpdateMenu()
	coin := ct.HighlightedRowCoin()
	if coin == nil {
		return nil
	}

	// read input field
	b := make([]byte, 100)
	n, err := ct.Views.Input.Read(b)
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}

	input := string(b[:n])
	holdings, err := eval.EvaluateExpressionToFloat64(input, coin)
	if err != nil {
		// leave value as is if expression can't be evaluated
		return nil
	}
	shouldDelete := holdings == 0

	idx := ct.GetPortfolioCoinIndex(coin)
	if err := ct.SetPortfolioEntry(coin.Name, holdings); err != nil {
		return err
	}

	if shouldDelete {
		if err := ct.RemovePortfolioEntry(coin.Name); err != nil {
			return err
		}
		ct.UpdateTable()
		if idx > 0 {
			idx -= 1
		}
	} else {
		ct.UpdateTable()
		ct.ToggleShowPortfolio()
		idx = ct.GetPortfolioCoinIndex(coin)
	}

	ct.HighlightRow(idx)

	return nil
}

// PortfolioEntry returns a portfolio entry
func (ct *Cointop) PortfolioEntry(c *Coin) (*PortfolioEntry, bool) {
	//log.Debug("PortfolioEntry()") // too many
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
	log.Debug("SetPortfolioEntry()")
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
func (ct *Cointop) RemovePortfolioEntry(coin string) error {
	log.Debug("RemovePortfolioEntry()")
	delete(ct.State.portfolio.Entries, strings.ToLower(coin))
	if err := ct.Save(); err != nil {
		return err
	}
	return nil
}

// PortfolioEntryExists returns true if portfolio entry exists
func (ct *Cointop) PortfolioEntryExists(c *Coin) bool {
	log.Debug("PortfolioEntryExists()")
	_, isNew := ct.PortfolioEntry(c)
	return !isNew
}

// PortfolioEntriesCount returns the count of portfolio entries
func (ct *Cointop) PortfolioEntriesCount() int {
	log.Debug("PortfolioEntriesCount()")
	return len(ct.State.portfolio.Entries)
}

// GetPortfolioSlice returns portfolio entries as a slice
func (ct *Cointop) GetPortfolioSlice() []*Coin {
	log.Debug("GetPortfolioSlice()")
	var sliced []*Coin
	if ct.PortfolioEntriesCount() == 0 {
		return sliced
	}

OUTER:
	for i := range ct.State.allCoins {
		coin := ct.State.allCoins[i]
		p, isNew := ct.PortfolioEntry(coin)
		if isNew {
			continue
		}
		// check not already found
		updateSlice := -1
		for j := range sliced {
			if coin.Symbol == sliced[j].Symbol {
				if coin.Rank >= sliced[j].Rank {
					continue OUTER // skip updates from lower-ranked coins
				}
				updateSlice = j // update this later
				break
			}
		}

		coin.Holdings = p.Holdings
		balance := coin.Price * p.Holdings
		balancestr := fmt.Sprintf("%.2f", balance)
		if ct.State.currencyConversion == "ETH" || ct.State.currencyConversion == "BTC" {
			balancestr = fmt.Sprintf("%.5f", balance)
		}
		balance, _ = strconv.ParseFloat(balancestr, 64)
		coin.Balance = balance
		if updateSlice == -1 {
			sliced = append(sliced, coin)
		} else {
			sliced[updateSlice] = coin
		}

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
	log.Debug("GetPortfolioTotal()")
	portfolio := ct.GetPortfolioSlice()
	var total float64
	for _, p := range portfolio {
		total += p.Balance
	}
	return total
}

// RefreshPortfolioCoins refreshes portfolio entry coin data
func (ct *Cointop) RefreshPortfolioCoins() error {
	log.Debug("RefreshPortfolioCoins()")
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
	SortBy           string
	SortDesc         bool
	HumanReadable    bool
	Format           string
	Filter           []string
	Cols             []string
	Convert          string
	NoHeader         bool
	PercentChange24H bool
	HideBalances     bool
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
	log.Debug("PrintHoldingsTable()")
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
	filterCoins := options.Filter
	filterCols := options.Cols
	holdings := ct.GetPortfolioSlice()
	noHeader := options.NoHeader
	hideBalances := options.HideBalances

	if format == "" {
		format = "table"
	}

	if sortBy != "" {
		if _, ok := portfolioColumns[sortBy]; !ok {
			return fmt.Errorf("the option %q is not a valid column name", sortBy)
		}

		ct.Sort(sortBy, sortDesc, holdings, true)
	}

	if _, ok := outputFormats[format]; !ok {
		return fmt.Errorf("the option %q is not a valid format type", format)
	}

	total := ct.GetPortfolioTotal()
	records := make([][]string, len(holdings))
	symbol := ct.CurrencySymbol()

	headers := []string{"name", "symbol", "price", "holdings", "balance", "24h%", "%holdings"}
	if len(filterCols) > 0 {
		for _, col := range filterCols {
			valid := false
			for _, h := range headers {
				if col == h {
					valid = true
					break
				}
			}
			switch col {
			case "amount":
				return fmt.Errorf("did you mean %q?", "balance")
			case "24H":
				fallthrough
			case "24H%":
				fallthrough
			case "24h":
				fallthrough
			case "24h_change":
				return fmt.Errorf("did you mean %q?", "24h%")
			case "percent_holdings":
				return fmt.Errorf("did you mean %q?", "%holdings")
			}
			if !valid {
				return fmt.Errorf("unsupported column value %q", col)
			}
		}
		headers = filterCols
	}

	for i, entry := range holdings {
		if len(filterCoins) > 0 {
			found := false
			for _, item := range filterCoins {
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

		item := make([]string, len(headers))
		for i, header := range headers {
			switch header {
			case "name":
				item[i] = entry.Name
			case "symbol":
				item[i] = entry.Symbol
			case "price":
				if humanReadable {
					item[i] = fmt.Sprintf("%s%s", symbol, ct.FormatPrice(entry.Price))
				} else {
					item[i] = strconv.FormatFloat(entry.Price, 'f', -1, 64)
				}
			case "holdings":
				if humanReadable {
					item[i] = humanize.Monetaryf(entry.Holdings, 2)
				} else {
					item[i] = strconv.FormatFloat(entry.Holdings, 'f', -1, 64)
				}
				if hideBalances {
					item[i] = HiddenBalanceChars
				}
			case "balance":
				if humanReadable {
					item[i] = fmt.Sprintf("%s%s", symbol, humanize.Monetaryf(entry.Balance, 2))
				} else {
					item[i] = strconv.FormatFloat(entry.Balance, 'f', -1, 64)
				}
				if hideBalances {
					item[i] = HiddenBalanceChars
				}
			case "24h%":
				if humanReadable {
					item[i] = fmt.Sprintf("%s%%", humanize.Numericf(entry.PercentChange24H, 2))
				} else {
					item[i] = fmt.Sprintf("%.2f", entry.PercentChange24H)
				}
			case "%holdings":
				if humanReadable {
					item[i] = fmt.Sprintf("%s%%", humanize.Numericf(percentHoldings, 2))
				} else {
					item[i] = fmt.Sprintf("%.2f", percentHoldings)
				}
				if hideBalances {
					item[i] = HiddenBalanceChars
				}
			}
		}
		records[i] = item
	}

	if format == "csv" {
		csvWriter := csv.NewWriter(os.Stdout)
		if !noHeader {
			if err := csvWriter.Write(headers); err != nil {
				return err
			}
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
		var output []byte
		var err error
		if noHeader {
			output, err = json.Marshal(records)
			if err != nil {
				return err
			}
		} else {
			list := make([]map[string]string, len(records))
			for i, record := range records {
				obj := make(map[string]string, len(record))
				for j, column := range record {
					obj[headers[j]] = column
				}

				list[i] = obj
			}

			output, err = json.Marshal(list)
			if err != nil {
				return err
			}
		}

		fmt.Println(string(output))
		return nil
	}

	alignment := []int{-1, -1, 1, 1, 1, 1, 1}
	var tableHeaders []string
	if !noHeader {
		tableHeaders = headers
	}
	table := asciitable.NewAsciiTable(&asciitable.Input{
		Data:      records,
		Headers:   tableHeaders,
		Alignment: alignment,
	})

	fmt.Println(table.String())
	return nil
}

// PrintHoldingsTotal prints the total holdings amount
func (ct *Cointop) PrintHoldingsTotal(options *TablePrintOptions) error {
	log.Debug("PrintHoldingsTotal()")
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
		if len(filter) > 0 {
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
		value = fmt.Sprintf("%s%s", symbol, humanize.Monetaryf(total, 2))
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

// PrintHoldings24HChange prints the total holdings amount
func (ct *Cointop) PrintHoldings24HChange(options *TablePrintOptions) error {
	log.Debug("PrintHoldings24HChange()")
	if options == nil {
		options = &TablePrintOptions{}
	}

	if err := ct.SetCurrencyConverstion(options.Convert); err != nil {
		return err
	}

	ct.RefreshPortfolioCoins()

	humanReadable := options.HumanReadable
	format := options.Format
	filter := options.Filter
	portfolio := ct.GetPortfolioSlice()
	total := ct.GetPortfolioTotal()
	var percentChange24H float64
	for _, entry := range portfolio {
		if len(filter) > 0 {
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

		n := (entry.Balance / total) * entry.PercentChange24H
		if math.IsNaN(n) {
			continue
		}
		percentChange24H += n
	}

	value := fmt.Sprintf("%.2f", percentChange24H)
	if humanReadable {
		value = fmt.Sprintf("%s%%", value)
	}

	if format == "csv" {
		csvWriter := csv.NewWriter(os.Stdout)
		if err := csvWriter.Write([]string{"24H%"}); err != nil {
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
		obj := map[string]string{"24H%": value}
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

// GetPortfolioCoinIndex returns the row index of coin in portfolio
func (ct *Cointop) GetPortfolioCoinIndex(coin *Coin) int {
	coins := ct.GetPortfolioSlice()
	for i, c := range coins {
		if c.ID == coin.ID {
			return i
		}
	}
	return 0
}

func (ct *Cointop) GetLastPortfolioRowIndex() int {
	l := ct.PortfolioLen()
	if l > 0 {
		l -= 1
	}
	return l
}

// IsPortfolioVisible returns true if portfolio view is visible
func (ct *Cointop) IsPortfolioVisible() bool {
	return ct.State.selectedView == PortfolioView
}

// PortfolioLen returns the number of portfolio entries
func (ct *Cointop) PortfolioLen() int {
	return len(ct.GetPortfolioSlice())
}

// TogglePortfolioBalances toggles hide/show portfolio balances. Useful for keeping balances secret when sharing screen or taking screenshots.
func (ct *Cointop) TogglePortfolioBalances() error {
	ct.State.hidePortfolioBalances = !ct.State.hidePortfolioBalances
	ct.UpdateUI(func() error {
		go ct.UpdateChart()
		go ct.UpdateTable()
		return nil
	})
	return nil
}
