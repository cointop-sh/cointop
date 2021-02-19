package cointop

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/miguelmota/cointop/pkg/pathutil"
	"github.com/miguelmota/cointop/pkg/toml"
)

var fileperm = os.FileMode(0644)

// ErrInvalidPriceAlert is error for invalid price alert value
var ErrInvalidPriceAlert = errors.New("invalid price alert value")

// NOTE: this is to support previous default config filepaths
var possibleConfigPaths = []string{
	":PREFERRED_CONFIG_HOME:/cointop/config.toml",
	":HOME:/.config/cointop/config.toml",
	":HOME:/.config/cointop/config",
	":HOME:/.cointop/config",
	":HOME:/.cointop/config.toml",
}

type config struct {
	Shortcuts     map[string]interface{} `toml:"shortcuts"`
	Favorites     map[string]interface{} `toml:"favorites"`
	Portfolio     map[string]interface{} `toml:"portfolio"`
	PriceAlerts   map[string]interface{} `toml:"price_alerts"`
	Currency      interface{}            `toml:"currency"`
	DefaultView   interface{}            `toml:"default_view"`
	CoinMarketCap map[string]interface{} `toml:"coinmarketcap"`
	API           interface{}            `toml:"api"`
	Colorscheme   interface{}            `toml:"colorscheme"`
	RefreshRate   interface{}            `toml:"refresh_rate"`
	CacheDir      interface{}            `toml:"cache_dir"`
	Table         map[string]interface{} `toml:"table"`
}

// SetupConfig loads config file
func (ct *Cointop) SetupConfig() error {
	ct.debuglog("setupConfig()")
	if err := ct.CreateConfigIfNotExists(); err != nil {
		return err
	}
	if err := ct.parseConfig(); err != nil {
		return err
	}
	if err := ct.loadTableColumnsFromConfig(); err != nil {
		return err
	}
	if err := ct.loadShortcutsFromConfig(); err != nil {
		return err
	}
	if err := ct.loadFavoritesFromConfig(); err != nil {
		return err
	}
	if err := ct.loadCurrencyFromConfig(); err != nil {
		return err
	}
	if err := ct.loadDefaultViewFromConfig(); err != nil {
		return err
	}
	if err := ct.loadAPIKeysFromConfig(); err != nil {
		return err
	}
	if err := ct.loadAPIChoiceFromConfig(); err != nil {
		return err
	}
	if err := ct.loadColorschemeFromConfig(); err != nil {
		return err
	}
	if err := ct.loadRefreshRateFromConfig(); err != nil {
		return err
	}
	if err := ct.loadCacheDirFromConfig(); err != nil {
		return err
	}
	if err := ct.loadPriceAlertsFromConfig(); err != nil {
		return err
	}
	if err := ct.loadPortfolioFromConfig(); err != nil {
		return err
	}

	return nil
}

// CreateConfigIfNotExists creates config file if it doesn't exist
func (ct *Cointop) CreateConfigIfNotExists() error {
	ct.debuglog("createConfigIfNotExists()")

	for _, configPath := range possibleConfigPaths {
		normalizedPath := pathutil.NormalizePath(configPath)
		if _, err := os.Stat(normalizedPath); err == nil {
			ct.configFilepath = normalizedPath
			return nil
		}
	}

	err := ct.makeConfigDir()
	if err != nil {
		return err
	}

	err = ct.makeConfigFile()
	if err != nil {
		return err
	}

	return nil
}

// ConfigDirPath returns the config directory path
func (ct *Cointop) ConfigDirPath() string {
	ct.debuglog("configDirPath()")
	path := pathutil.NormalizePath(ct.configFilepath)
	separator := string(filepath.Separator)
	parts := strings.Split(path, separator)
	return strings.Join(parts[0:len(parts)-1], separator)
}

// ConfigFilePath return the config file path
func (ct *Cointop) ConfigFilePath() string {
	ct.debuglog("configFilePath()")
	return pathutil.NormalizePath(ct.configFilepath)
}

// ConfigPath return the config file path
func (ct *Cointop) makeConfigDir() error {
	ct.debuglog("makeConfigDir()")
	path := ct.ConfigDirPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

// MakeConfigFile creates a new config file
func (ct *Cointop) makeConfigFile() error {
	ct.debuglog("makeConfigFile()")
	path := ct.ConfigFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fo, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fo.Close()
		b, err := ct.configToToml()
		if err != nil {
			return err
		}
		if _, err := fo.Write(b); err != nil {
			return err
		}
	}
	return nil
}

// SaveConfig writes settings to the config file
func (ct *Cointop) SaveConfig() error {
	ct.debuglog("saveConfig()")
	ct.saveMux.Lock()
	defer ct.saveMux.Unlock()
	path := ct.ConfigFilePath()
	if _, err := os.Stat(path); err == nil {
		b, err := ct.configToToml()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, b, fileperm)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseConfig decodes the toml config file
func (ct *Cointop) parseConfig() error {
	ct.debuglog("parseConfig()")
	var conf config
	path := ct.ConfigFilePath()
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return err
	}

	ct.config = conf
	return nil
}

// ConfigToToml encodes config struct to TOML
func (ct *Cointop) configToToml() ([]byte, error) {
	ct.debuglog("configToToml()")
	shortcutsIfcs := map[string]interface{}{}
	for k, v := range ct.State.shortcutKeys {
		var i interface{} = v
		shortcutsIfcs[k] = i
	}

	var favoritesIfc []interface{}
	for k, ok := range ct.State.favorites {
		if ok {
			var i interface{} = k
			favoritesIfc = append(favoritesIfc, i)
		}
	}
	sort.Slice(favoritesIfc, func(i, j int) bool {
		return favoritesIfc[i].(string) < favoritesIfc[j].(string)
	})

	var favoritesBySymbolIfc []interface{}
	favoritesMapIfc := map[string]interface{}{
		// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
		"symbols": favoritesBySymbolIfc,
		"names":   favoritesIfc,
	}

	var favoritesColumnsIfc interface{} = ct.State.favoritesTableColumns
	favoritesMapIfc["columns"] = favoritesColumnsIfc

	portfolioIfc := map[string]interface{}{}
	var holdingsIfc [][]string
	for name := range ct.State.portfolio.Entries {
		entry, ok := ct.State.portfolio.Entries[name]
		if !ok || entry.Coin == "" {
			continue
		}
		var amount string = strconv.FormatFloat(entry.Holdings, 'f', -1, 64)
		var coinName string = entry.Coin
		var tuple []string = []string{coinName, amount}
		holdingsIfc = append(holdingsIfc, tuple)
	}
	sort.Slice(holdingsIfc, func(i, j int) bool {
		return holdingsIfc[i][0] < holdingsIfc[j][0]
	})
	portfolioIfc["holdings"] = holdingsIfc

	var columnsIfc interface{} = ct.State.portfolioTableColumns
	portfolioIfc["columns"] = columnsIfc

	var currencyIfc interface{} = ct.State.currencyConversion
	var defaultViewIfc interface{} = ct.State.defaultView
	var colorschemeIfc interface{} = ct.colorschemeName
	var refreshRateIfc interface{} = uint(ct.State.refreshRate.Seconds())
	var cacheDirIfc interface{} = ct.State.cacheDir

	cmcIfc := map[string]interface{}{
		"pro_api_key": ct.apiKeys.cmc,
	}

	var apiChoiceIfc interface{} = ct.apiChoice

	var priceAlertsIfc []interface{}
	for _, priceAlert := range ct.State.priceAlerts.Entries {
		if priceAlert.Expired {
			continue
		}
		priceAlertsIfc = append(priceAlertsIfc, []string{
			priceAlert.CoinName,
			priceAlert.Operator,
			strconv.FormatFloat(priceAlert.TargetPrice, 'f', -1, 64),
			priceAlert.Frequency,
		})
	}
	priceAlertsMapIfc := map[string]interface{}{
		"alerts": priceAlertsIfc,
		//"sound":  ct.State.priceAlerts.SoundEnabled,
	}

	var coinsTableColumnsIfc interface{} = ct.State.coinsTableColumns
	tableMapIfc := map[string]interface{}{}
	tableMapIfc["columns"] = coinsTableColumnsIfc

	var inputs = &config{
		API:           apiChoiceIfc,
		Colorscheme:   colorschemeIfc,
		CoinMarketCap: cmcIfc,
		Currency:      currencyIfc,
		DefaultView:   defaultViewIfc,
		Favorites:     favoritesMapIfc,
		RefreshRate:   refreshRateIfc,
		Shortcuts:     shortcutsIfcs,
		Portfolio:     portfolioIfc,
		PriceAlerts:   priceAlertsMapIfc,
		CacheDir:      cacheDirIfc,
		Table:         tableMapIfc,
	}

	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	err := encoder.Encode(inputs)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// LoadTableColumnsFromConfig loads preferred coins table columns from config file to struct
func (ct *Cointop) loadTableColumnsFromConfig() error {
	ct.debuglog("loadTableColumnsFromConfig()")
	columnsIfc, ok := ct.config.Table["columns"]
	if !ok {
		return nil
	}
	var columns []string
	ifcs, ok := columnsIfc.([]interface{})
	if ok {
		for _, ifc := range ifcs {
			if v, ok := ifc.(string); ok {
				if !ct.ValidCoinsTableHeader(v) {
					return fmt.Errorf("invalid table header name %q. Valid names are: %s", v, strings.Join(SupportedCoinTableHeaders, ","))
				}
				columns = append(columns, v)
			}
		}
		if len(columns) > 0 {
			ct.State.coinsTableColumns = columns
		}
	}
	return nil
}

// LoadShortcutsFromConfig loads keyboard shortcuts from config file to struct
func (ct *Cointop) loadShortcutsFromConfig() error {
	ct.debuglog("loadShortcutsFromConfig()")
	for k, ifc := range ct.config.Shortcuts {
		if v, ok := ifc.(string); ok {
			if !ct.ActionExists(v) {
				continue
			}
			ct.State.shortcutKeys[k] = v
		}
	}
	return nil
}

// LoadCurrencyFromConfig loads currency from config file to struct
func (ct *Cointop) loadCurrencyFromConfig() error {
	ct.debuglog("loadCurrencyFromConfig()")
	if currency, ok := ct.config.Currency.(string); ok {
		ct.State.currencyConversion = strings.ToUpper(currency)
	}
	return nil
}

// LoadDefaultViewFromConfig loads default view from config file to struct
func (ct *Cointop) loadDefaultViewFromConfig() error {
	ct.debuglog("loadDefaultViewFromConfig()")
	if defaultView, ok := ct.config.DefaultView.(string); ok {
		defaultView = strings.ToLower(defaultView)
		switch defaultView {
		case "portfolio":
			ct.SetSelectedView(PortfolioView)
		case "favorites":
			ct.SetSelectedView(FavoritesView)
		case "alerts", "price_alerts":
			ct.SetSelectedView(PriceAlertsView)
		case "default":
			fallthrough
		default:
			ct.SetSelectedView(CoinsView)
			defaultView = "default"
		}
		ct.State.defaultView = defaultView
	}
	return nil
}

// LoadAPIKeysFromConfig loads API keys from config file to struct
func (ct *Cointop) loadAPIKeysFromConfig() error {
	ct.debuglog("loadAPIKeysFromConfig()")
	for key, value := range ct.config.CoinMarketCap {
		k := strings.TrimSpace(strings.ToLower(key))
		if k == "pro_api_key" {
			ct.apiKeys.cmc = value.(string)
		}
	}
	return nil
}

// LoadColorschemeFromConfig loads colorscheme name from config file to struct
func (ct *Cointop) loadColorschemeFromConfig() error {
	ct.debuglog("loadColorschemeFromConfig()")
	if colorscheme, ok := ct.config.Colorscheme.(string); ok {
		ct.colorschemeName = colorscheme
	}

	return nil
}

// LoadRefreshRateFromConfig loads refresh rate from config file to struct
func (ct *Cointop) loadRefreshRateFromConfig() error {
	ct.debuglog("loadRefreshRateFromConfig()")
	if refreshRate, ok := ct.config.RefreshRate.(int64); ok {
		ct.State.refreshRate = time.Duration(uint(refreshRate)) * time.Second
	}

	return nil
}

// LoadCacheDirFromConfig loads cache dir from config file to struct
func (ct *Cointop) loadCacheDirFromConfig() error {
	ct.debuglog("loadCacheDirFromConfig()")
	if cacheDir, ok := ct.config.CacheDir.(string); ok {
		ct.State.cacheDir = pathutil.NormalizePath(cacheDir)
	}

	return nil
}

// GetColorschemeColors loads colors from colorsheme file to struct
func (ct *Cointop) getColorschemeColors() (map[string]interface{}, error) {
	ct.debuglog("getColorschemeColors()")
	var colors map[string]interface{}
	if ct.colorschemeName == "" {
		ct.colorschemeName = DefaultColorscheme
		if _, err := toml.Decode(DefaultColors, &colors); err != nil {
			return nil, err
		}
	} else {
		colorsDir := fmt.Sprintf("%s/colors", ct.ConfigDirPath())
		if ct.colorsDir != "" {
			colorsDir = pathutil.NormalizePath(ct.colorsDir)
		}

		path := fmt.Sprintf("%s/%s.toml", colorsDir, ct.colorschemeName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// NOTE: case for when cointop is set as the theme but the colorscheme file doesn't exist
			if ct.colorschemeName == "cointop" {
				if _, err := toml.Decode(DefaultColors, &colors); err != nil {
					return nil, err
				}

				return colors, nil
			}

			return nil, fmt.Errorf("the colorscheme file %q was not found.\n%s", path, ColorschemeHelpString())
		}

		if _, err := toml.DecodeFile(path, &colors); err != nil {
			return nil, err
		}
	}

	return colors, nil
}

// LoadAPIChoiceFromConfig loads API choices from config file to struct
func (ct *Cointop) loadAPIChoiceFromConfig() error {
	ct.debuglog("loadAPIKeysFromConfig()")
	apiChoice, ok := ct.config.API.(string)
	if ok {
		apiChoice = strings.TrimSpace(strings.ToLower(apiChoice))
		ct.apiChoice = apiChoice
	}
	return nil
}

// LoadFavoritesFromConfig loads favorites data from config file to struct
func (ct *Cointop) loadFavoritesFromConfig() error {
	ct.debuglog("loadFavoritesFromConfig()")
	for k, valueIfc := range ct.config.Favorites {
		ifcs, ok := valueIfc.([]interface{})
		if !ok {
			continue
		}
		switch k {
		// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
		case "symbols":
			for _, ifc := range ifcs {
				if v, ok := ifc.(string); ok {
					ct.State.favoritesBySymbol[strings.ToUpper(v)] = true
				}
			}
		case "names":
			for _, ifc := range ifcs {
				if v, ok := ifc.(string); ok {
					ct.State.favorites[v] = true
				}
			}
		case "columns":
			var columns []string
			for _, ifc := range ifcs {
				col, ok := ifc.(string)
				if !ok {
					continue
				}
				if !ct.ValidCoinsTableHeader(col) {
					return fmt.Errorf("invalid table header name %q. Valid names are: %s", col, strings.Join(SupportedCoinTableHeaders, ","))
				}
				columns = append(columns, col)
			}
			if len(columns) > 0 {
				ct.State.favoritesTableColumns = columns
			}
		}
	}
	return nil
}

// LoadPortfolioFromConfig loads portfolio data from config file to struct
func (ct *Cointop) loadPortfolioFromConfig() error {
	ct.debuglog("loadPortfolioFromConfig()")

	for key, valueIfc := range ct.config.Portfolio {
		if key == "columns" {
			var columns []string
			ifcs, ok := valueIfc.([]interface{})
			if ok {
				for _, ifc := range ifcs {
					if v, ok := ifc.(string); ok {
						if !ct.ValidPortfolioTableHeader(v) {
							return fmt.Errorf("invalid table header name %q. Valid names are: %s", v, strings.Join(SupportedPortfolioTableHeaders, ","))
						}
						columns = append(columns, v)
					}
				}
				if len(columns) > 0 {
					ct.State.portfolioTableColumns = columns
				}
			}
		} else if key == "holdings" {
			holdingsIfc, ok := valueIfc.([]interface{})
			if !ok {
				continue
			}

			for _, itemIfc := range holdingsIfc {
				tupleIfc, ok := itemIfc.([]interface{})
				if !ok {
					continue
				}
				if len(tupleIfc) > 2 {
					continue
				}
				name, ok := tupleIfc[0].(string)
				if !ok {
					continue
				}

				holdings, err := ct.InterfaceToFloat64(tupleIfc[1])
				if err != nil {
					return nil
				}

				if err := ct.SetPortfolioEntry(name, holdings); err != nil {
					return err
				}
			}
		} else {
			// Backward compatibility < v1.6.0
			holdings, err := ct.InterfaceToFloat64(valueIfc)
			if err != nil {
				return err
			}

			if err := ct.SetPortfolioEntry(key, holdings); err != nil {
				return err
			}
		}
	}

	return nil
}

// LoadPriceAlertsFromConfig loads price alerts from config file to struct
func (ct *Cointop) loadPriceAlertsFromConfig() error {
	ct.debuglog("loadPriceAlertsFromConfig()")
	priceAlertsIfc, ok := ct.config.PriceAlerts["alerts"]
	if !ok {
		return nil
	}
	priceAlertsSliceIfc, ok := priceAlertsIfc.([]interface{})
	if !ok {
		return nil
	}
	for _, priceAlertIfc := range priceAlertsSliceIfc {
		priceAlert, ok := priceAlertIfc.([]interface{})
		if !ok {
			return ErrInvalidPriceAlert
		}
		coinName, ok := priceAlert[0].(string)
		if !ok {
			return ErrInvalidPriceAlert
		}
		operator, ok := priceAlert[1].(string)
		if !ok {
			return ErrInvalidPriceAlert
		}
		if _, ok := PriceAlertOperatorMap[operator]; !ok {
			return ErrInvalidPriceAlert
		}
		targetPrice, err := ct.InterfaceToFloat64(priceAlert[2])
		if err != nil {
			return err
		}
		frequency, ok := priceAlert[3].(string)
		if !ok {
			return ErrInvalidPriceAlert
		}
		if _, ok := PriceAlertFrequencyMap[frequency]; !ok {
			return ErrInvalidPriceAlert
		}
		id := strings.ToLower(fmt.Sprintf("%s_%s_%v_%s", coinName, operator, targetPrice, frequency))
		entry := &PriceAlert{
			ID:          id,
			CoinName:    coinName,
			Operator:    operator,
			TargetPrice: targetPrice,
			Frequency:   frequency,
		}
		ct.State.priceAlerts.Entries = append(ct.State.priceAlerts.Entries, entry)
	}
	soundIfc, ok := ct.config.PriceAlerts["sound"]
	if ok {
		enabled, ok := soundIfc.(bool)
		if !ok {
			return ErrInvalidPriceAlert
		}
		ct.State.priceAlerts.SoundEnabled = enabled
	}

	return nil
}

// InterfaceToFloat64 attempts to convert interface to float64
func (ct *Cointop) InterfaceToFloat64(value interface{}) (float64, error) {
	var num float64
	var err error
	switch v := value.(type) {
	case string:
		num, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
	case int:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case float64:
		num = v
	}

	return num, nil
}
