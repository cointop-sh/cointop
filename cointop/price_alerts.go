package cointop

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/cointop/pkg/humanize"
	"github.com/miguelmota/cointop/pkg/notifier"
	"github.com/miguelmota/cointop/pkg/pad"
	"github.com/miguelmota/cointop/pkg/table"
)

// GetPriceAlertsTableHeaders returns the alerts table headers
func (ct *Cointop) GetPriceAlertsTableHeaders() []string {
	return []string{
		"name",
		"symbol",
		"target_price",
		"price",
		"frequency",
	}
}

// PriceAlertOperatorMap is map of valid price alert operator symbols
var PriceAlertOperatorMap = map[string]string{
	">":  ">",
	"<":  "<",
	">=": "≥",
	"<=": "≤",
	"=":  "=",
}

// PriceAlertFrequencyMap is map of valid price alert frequency values
var PriceAlertFrequencyMap = map[string]bool{
	"once":        true,
	"reoccurring": true,
}

// GetPriceAlertsTable returns the table for displaying alerts
func (ct *Cointop) GetPriceAlertsTable() *table.Table {
	ct.debuglog("getPriceAlertsTable()")
	maxX := ct.width()
	t := table.NewTable().SetWidth(maxX)
	var rows [][]*table.RowCell
	headers := ct.GetPriceAlertsTableHeaders()
	ct.ClearSyncMap(ct.State.tableColumnWidths)
	ct.ClearSyncMap(ct.State.tableColumnAlignLeft)
	for _, entry := range ct.State.priceAlerts.Entries {
		if entry.Expired {
			continue
		}
		ifc, ok := ct.State.allCoinsSlugMap.Load(entry.CoinName)
		if !ok {
			continue
		}
		coin, ok := ifc.(*Coin)
		if !ok {
			continue
		}
		_, ok = PriceAlertOperatorMap[entry.Operator]
		if !ok {
			continue
		}

		leftMargin := 1
		rightMargin := 1
		var rowCells []*table.RowCell
		for _, header := range headers {
			switch header {
			case "name":
				name := TruncateString(entry.CoinName, 16)
				ct.SetTableColumnWidthFromString(header, name)
				ct.SetTableColumnAlignLeft(header, true)
				namecolor := ct.colorscheme.TableRow
				rowCells = append(rowCells, &table.RowCell{
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
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   true,
					Color:       ct.colorscheme.TableRow,
					Text:        symbol,
				})

			case "target_price":
				targetPrice := fmt.Sprintf("%s %s", entry.Operator, humanize.Commaf(entry.TargetPrice))
				ct.SetTableColumnWidthFromString(header, targetPrice)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       ct.colorscheme.TableColumnPrice,
					Text:        targetPrice,
				})
			case "price":
				text := humanize.Commaf(coin.Price)
				ct.SetTableColumnWidthFromString(header, text)
				ct.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       ct.colorscheme.TableRow,
					Text:        text,
				})
			case "frequency":
				frequency := entry.Frequency
				ct.SetTableColumnWidthFromString(header, frequency)
				ct.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells, &table.RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   true,
					Color:       ct.colorscheme.TableRow,
					Text:        frequency,
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

// TogglePriceAlerts toggles the price alerts view
func (ct *Cointop) TogglePriceAlerts() error {
	ct.debuglog("togglePriceAlerts()")
	ct.ToggleSelectedView(PriceAlertsView)
	ct.NavigateFirstLine()
	go ct.UpdateTable()
	return nil
}

// IsPriceAlertsVisible returns true if alerts view is visible
func (ct *Cointop) IsPriceAlertsVisible() bool {
	return ct.State.selectedView == PriceAlertsView
}

// PriceAlertWatcher starts the price alert watcher
func (ct *Cointop) PriceAlertWatcher() error {
	ct.debuglog("priceAlertWatcher()")
	alerts := ct.State.priceAlerts.Entries
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		for _, alert := range alerts {
			err := ct.CheckPriceAlert(alert)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckPriceAlert checks the price alert
func (ct *Cointop) CheckPriceAlert(alert *PriceAlert) error {
	ct.debuglog("checkPriceAlert()")
	if alert.Expired {
		return nil
	}

	coinIfc, _ := ct.State.allCoinsSlugMap.Load(alert.CoinName)
	coin, ok := coinIfc.(*Coin)
	if !ok {
		return nil
	}
	var msg string
	title := "Cointop Alert"
	priceStr := fmt.Sprintf("%s%s (%s%s)", ct.CurrencySymbol(), humanize.Commaf(alert.TargetPrice), ct.CurrencySymbol(), humanize.Commaf(coin.Price))
	if alert.Operator == ">" {
		if coin.Price > alert.TargetPrice {
			msg = fmt.Sprintf("%s price is greater than %v", alert.CoinName, priceStr)
		}
	} else if alert.Operator == ">=" {
		if coin.Price >= alert.TargetPrice {
			msg = fmt.Sprintf("%s price is greater than or equal to %v", alert.CoinName, priceStr)
		}
	} else if alert.Operator == "<" {
		if coin.Price < alert.TargetPrice {
			msg = fmt.Sprintf("%s price is less than %v", alert.CoinName, priceStr)
		}
	} else if alert.Operator == "<=" {
		if coin.Price <= alert.TargetPrice {
			msg = fmt.Sprintf("%s price is less than or equal to %v", alert.CoinName, priceStr)
		}
	} else if alert.Operator == "=" {
		if coin.Price == alert.TargetPrice {
			msg = fmt.Sprintf("%s price is equal to %v", alert.CoinName, priceStr)
		}
	}

	if msg != "" {
		if ct.State.priceAlerts.SoundEnabled {
			notifier.Notify(title, msg)
		} else {
			notifier.Notify(title, msg)
		}

		alert.Expired = true
	}

	if err := ct.Save(); err != nil {
		return err
	}
	return nil
}

// UpdatePriceAlertsUpdateMenu updates the alerts update menu view
func (ct *Cointop) UpdatePriceAlertsUpdateMenu(isNew bool) error {
	ct.debuglog("updatePriceAlertsUpdateMenu()")

	exists := false
	var value string
	var currentPrice string
	var coinName string
	ct.State.priceAlertEditID = ""
	if !isNew && ct.IsPriceAlertsVisible() {
		rowIndex := ct.HighlightedRowIndex()
		entry := ct.State.priceAlerts.Entries[rowIndex]
		ifc, ok := ct.State.allCoinsSlugMap.Load(entry.CoinName)
		if ok {
			coin, ok := ifc.(*Coin)
			if ok {
				coinName = entry.CoinName
				currentPrice = strconv.FormatFloat(coin.Price, 'f', -1, 64)
				value = fmt.Sprintf("%s %v", entry.Operator, entry.TargetPrice)
				ct.State.priceAlertEditID = entry.ID
				exists = true
			}
		}
	}

	var mode string
	var current string
	var submitText string
	var offset int
	if exists {
		mode = "Edit"
		current = fmt.Sprintf("(current %s%s)", ct.CurrencySymbol(), currentPrice)
		submitText = "Set"
		offset = ct.width() - 21
	} else {
		coin := ct.HighlightedRowCoin()
		coinName = coin.Name
		currentPrice = strconv.FormatFloat(coin.Price, 'f', -1, 64)
		value = fmt.Sprintf("> %s", currentPrice)
		mode = "Create"
		submitText = "Create"
		offset = ct.width() - 23
	}
	header := ct.colorscheme.MenuHeader(fmt.Sprintf(" %s Alert Entry %s\n\n", mode, pad.Left("[q] close ", offset, " ")))
	label := fmt.Sprintf(" Enter target price for %s %s", ct.colorscheme.MenuLabel(coinName), current)
	content := fmt.Sprintf("%s\n%s\n\n%s%s\n\n\n [Enter] %s    [ESC] Cancel", header, label, strings.Repeat(" ", 29), ct.State.currencyConversion, submitText)

	ct.UpdateUI(func() error {
		ct.Views.Menu.SetFrame(true)
		ct.Views.Menu.Update(content)
		ct.Views.Input.Write(value)
		ct.Views.Input.SetCursor(len(value), 0)
		return nil
	})
	return nil
}

// ShowPriceAlertsAddMenu shows the alert add menu
func (ct *Cointop) ShowPriceAlertsAddMenu() error {
	ct.debuglog("showPriceAlertsAddMenu()")
	ct.SetSelectedView(PriceAlertsView)
	ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	ct.UpdatePriceAlertsUpdateMenu(true)
	ct.ui.SetCursor(true)
	ct.SetActiveView(ct.Views.Menu.Name())
	ct.g.SetViewOnTop(ct.Views.Input.Name())
	ct.g.SetCurrentView(ct.Views.Input.Name())
	return nil
}

// ShowPriceAlertsUpdateMenu shows the alerts update menu
func (ct *Cointop) ShowPriceAlertsUpdateMenu() error {
	ct.debuglog("showPriceAlertsUpdateMenu()")
	ct.SetSelectedView(PriceAlertsView)
	ct.State.lastSelectedRowIndex = ct.HighlightedPageRowIndex()
	ct.UpdatePriceAlertsUpdateMenu(false)
	ct.ui.SetCursor(true)
	ct.SetActiveView(ct.Views.Menu.Name())
	ct.g.SetViewOnTop(ct.Views.Input.Name())
	ct.g.SetCurrentView(ct.Views.Input.Name())
	return nil
}

// HidePriceAlertsUpdateMenu hides the alerts update menu
func (ct *Cointop) HidePriceAlertsUpdateMenu() error {
	ct.debuglog("hidePriceAlertsUpdateMenu()")
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

// EnterKeyPressHandler is the key press handle for update menus
func (ct *Cointop) EnterKeyPressHandler() error {
	if ct.IsPriceAlertsVisible() {
		return ct.CreatePriceAlert()
	}

	return ct.SetPortfolioHoldings()
}

// CreatePriceAlert sets price from inputed value
func (ct *Cointop) CreatePriceAlert() error {
	ct.debuglog("createPriceAlert()")
	defer ct.HidePriceAlertsUpdateMenu()
	var coinName string

	isNew := ct.State.priceAlertEditID == ""
	if isNew {
		coin := ct.HighlightedRowCoin()
		coinName = coin.Name
	} else {
		for i, entry := range ct.State.priceAlerts.Entries {
			if entry.ID == ct.State.priceAlertEditID {
				coinName = ct.State.priceAlerts.Entries[i].CoinName
			}
		}
	}

	operator, targetPrice, err := ct.ReadAndParsePriceAlertInput()
	if err != nil {
		return err
	}

	if err := ct.SetPriceAlert(coinName, operator, targetPrice); err != nil {
		return err
	}

	ct.UpdateTable()
	if isNew {
		ct.GoToPageRowIndex(0)
	}

	return nil
}

// ReadAndParsePriceAlertInput reads and parses price alert input field value
func (ct *Cointop) ReadAndParsePriceAlertInput() (string, float64, error) {
	// read input field
	b := make([]byte, 100)
	n, err := ct.Views.Input.Read(b)
	if err != nil {
		return "", 0, err
	}
	if n == 0 {
		return "", 0, nil
	}

	inputValue := string(b)
	operator, targetPrice, err := ct.ParsePriceAlertInput(inputValue)
	if err != nil {
		return "", 0, err
	}

	return operator, targetPrice, nil
}

// ParsePriceAlertInput parses price alert input field value
func (ct *Cointop) ParsePriceAlertInput(value string) (string, float64, error) {
	regex := regexp.MustCompile(`(>|<|>=|<=|=)?\s*([0-9.]+).*`)
	matches := regex.FindStringSubmatch(strings.TrimSpace(value))
	operator := ""
	amountValue := ""
	if len(matches) == 2 {
		amountValue = matches[1]
	} else if len(matches) == 3 {
		operator = matches[1]
		amountValue = matches[2]
	}
	amountValue = normalizeFloatString(amountValue)
	targetPrice, err := strconv.ParseFloat(amountValue, 64)
	if err != nil {
		return "", 0, err
	}

	return operator, targetPrice, nil
}

// SetPriceAlert sets a price alert
func (ct *Cointop) SetPriceAlert(coinName string, operator string, targetPrice float64) error {
	ct.debuglog("setPriceAlert()")

	if operator == "" {
		operator = "="
	}

	if _, ok := PriceAlertOperatorMap[operator]; !ok {
		return errors.New("price alert operator is invalid")
	}

	frequency := "once"
	id := strings.ToLower(fmt.Sprintf("%s_%s_%v_%s", coinName, operator, targetPrice, frequency))
	newEntry := &PriceAlert{
		ID:          id,
		CoinName:    coinName,
		Operator:    operator,
		TargetPrice: targetPrice,
		Frequency:   frequency,
	}

	if ct.State.priceAlertEditID == "" {
		ct.State.priceAlerts.Entries = append([]*PriceAlert{newEntry}, ct.State.priceAlerts.Entries...)
	} else {
		for i, entry := range ct.State.priceAlerts.Entries {
			if entry.ID == ct.State.priceAlertEditID {
				ct.State.priceAlerts.Entries[i] = newEntry
			}
		}
	}

	if err := ct.Save(); err != nil {
		return err
	}

	return nil
}

// ActivePriceAlerts returns the active price alerts
func (ct *Cointop) ActivePriceAlerts() []*PriceAlert {
	var filtered []*PriceAlert
	for _, entry := range ct.State.priceAlerts.Entries {
		if entry.Expired {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

// ActivePriceAlertsLen returns the number of active price alerts
func (ct *Cointop) ActivePriceAlertsLen() int {
	return len(ct.ActivePriceAlerts())
}
