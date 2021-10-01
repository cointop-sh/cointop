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
	"unicode/utf8"

	"github.com/cointop-sh/cointop/pkg/pathutil"
	"github.com/cointop-sh/cointop/pkg/toml"
	log "github.com/sirupsen/logrus"
)

// FilePerm is the default file permissions
var FilePerm = os.FileMode(0644)

// ErrInvalidPriceAlert is error for invalid price alert value
var ErrInvalidPriceAlert = errors.New("invalid price alert value")

// PossibleConfigPaths are the the possible config file paths.
// NOTE: this is to support previous default config filepaths
var PossibleConfigPaths = []string{
	":PREFERRED_CONFIG_HOME:/cointop/config.toml",
	":HOME:/.config/cointop/config.toml",
	":HOME:/.config/cointop/config",
	":HOME:/.cointop/config",
	":HOME:/.cointop/config.toml",
}

// ConfigFileConfig is the config file structure
type ConfigFileConfig struct {
	Shortcuts         map[string]interface{} `toml:"shortcuts"`
	Favorites         map[string]interface{} `toml:"favorites"`
	Portfolio         map[string]interface{} `toml:"portfolio"`
	PriceAlerts       map[string]interface{} `toml:"price_alerts"`
	Currency          interface{}            `toml:"currency"`
	DefaultView       interface{}            `toml:"default_view"`
	DefaultChartRange interface{}            `toml:"default_chart_range"`
	CoinMarketCap     map[string]interface{} `toml:"coinmarketcap"`
	API               interface{}            `toml:"api"`
	Colorscheme       interface{}            `toml:"colorscheme"`
	RefreshRate       interface{}            `toml:"refresh_rate"`
	CacheDir          interface{}            `toml:"cache_dir"`
	Table             map[string]interface{} `toml:"table"`
	Chart             map[string]interface{} `toml:"chart"`
}

// SetupConfig loads config file
func (ct *Cointop) SetupConfig() error {
	type loadConfigFunc func() error
	loaders := []loadConfigFunc{
		ct.CreateConfigIfNotExists,
		ct.ParseConfig,
		ct.loadTableConfig,
		ct.loadChartConfig,
		ct.loadShortcutsFromConfig,
		ct.loadFavoritesFromConfig,
		ct.loadCurrencyFromConfig,
		ct.loadDefaultViewFromConfig,
		ct.loadDefaultChartRangeFromConfig,
		ct.loadAPIKeysFromConfig,
		ct.loadAPIChoiceFromConfig,
		ct.loadColorschemeFromConfig,
		ct.loadRefreshRateFromConfig,
		ct.loadCacheDirFromConfig,
		ct.loadPriceAlertsFromConfig,
		ct.loadPortfolioFromConfig,
	}

	for _, f := range loaders {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

// CreateConfigIfNotExists creates config file if it doesn't exist
func (ct *Cointop) CreateConfigIfNotExists() error {
	log.Debug("CreateConfigIfNotExists()")

	ct.configFilepath = pathutil.NormalizePath(ct.configFilepath)

	// check if config file exists in one of th default paths
	if ct.configFilepath == DefaultConfigFilepath {
		for _, configPath := range PossibleConfigPaths {
			normalizedPath := pathutil.NormalizePath(configPath)
			if _, err := os.Stat(normalizedPath); err == nil {
				ct.configFilepath = normalizedPath
				return nil
			}
		}
	}

	err := ct.MakeConfigDir()
	if err != nil {
		return err
	}

	err = ct.MakeConfigFile()
	if err != nil {
		return err
	}

	return nil
}

// ConfigDirPath returns the config directory path
func (ct *Cointop) ConfigDirPath() string {
	log.Debug("ConfigDirPath()")
	path := pathutil.NormalizePath(ct.configFilepath)
	separator := string(filepath.Separator)
	parts := strings.Split(path, separator)
	return strings.Join(parts[0:len(parts)-1], separator)
}

// ConfigFilePath return the config file path
func (ct *Cointop) ConfigFilePath() string {
	log.Debug("ConfigFilePath()")
	return pathutil.NormalizePath(ct.configFilepath)
}

// ConfigPath return the config file path
func (ct *Cointop) MakeConfigDir() error {
	log.Debug("MakeConfigDir()")
	path := ct.ConfigDirPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}

	return nil
}

// MakeConfigFile creates a new config file
func (ct *Cointop) MakeConfigFile() error {
	log.Debug("MakeConfigFile()")
	path := ct.ConfigFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fo, err := os.Create(path)
		if err != nil {
			return err
		}
		defer fo.Close()
		b, err := ct.ConfigToToml()
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
	log.Debug("SaveConfig()")
	ct.saveMux.Lock()
	defer ct.saveMux.Unlock()
	path := ct.ConfigFilePath()
	if _, err := os.Stat(path); err == nil {
		b, err := ct.ConfigToToml()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, b, FilePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseConfig decodes the toml config file
func (ct *Cointop) ParseConfig() error {
	log.Debug("ParseConfig()")
	var conf ConfigFileConfig
	path := ct.configFilepath
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return err
	}

	ct.config = conf
	return nil
}

// ConfigToToml encodes config struct to TOML
func (ct *Cointop) ConfigToToml() ([]byte, error) {
	log.Debug("ConfigToToml()")
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
		"columns": ct.State.favoritesTableColumns,
	}

	var favoritesColumnsIfc interface{} = ct.State.favoritesTableColumns
	favoritesMapIfc["columns"] = favoritesColumnsIfc
	favoritesMapIfc["character"] = ct.State.favoriteChar

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
	portfolioIfc := map[string]interface{}{
		"holdings": holdingsIfc,
		"columns":  ct.State.portfolioTableColumns,
	}

	cmcIfc := map[string]interface{}{
		"pro_api_key": ct.apiKeys.cmc,
	}

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

	tableMapIfc := map[string]interface{}{
		"columns":                ct.State.coinsTableColumns,
		"keep_row_focus_on_sort": ct.State.keepRowFocusOnSort,
	}

	chartMapIfc := map[string]interface{}{
		"max_width": ct.State.maxChartWidth,
		"height":    ct.State.chartHeight,
	}

	var inputs = &ConfigFileConfig{
		API:               ct.apiChoice,
		Colorscheme:       ct.colorschemeName,
		CoinMarketCap:     cmcIfc,
		Currency:          ct.State.currencyConversion,
		DefaultView:       ct.State.defaultView,
		DefaultChartRange: ct.State.defaultChartRange,
		Favorites:         favoritesMapIfc,
		RefreshRate:       uint(ct.State.refreshRate.Seconds()),
		Shortcuts:         shortcutsIfcs,
		Portfolio:         portfolioIfc,
		PriceAlerts:       priceAlertsMapIfc,
		CacheDir:          ct.State.cacheDir,
		Table:             tableMapIfc,
		Chart:             chartMapIfc,
	}

	var b bytes.Buffer
	encoder := toml.NewEncoder(&b)
	err := encoder.Encode(inputs)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// LoadTableConfig loads table config from toml config into state struct
func (ct *Cointop) loadTableConfig() error {
	log.Debug("loadTableConfig()")
	err := ct.loadTableColumnsFromConfig()
	if err != nil {
		return err
	}

	keepRowFocusOnSortIfc, ok := ct.config.Table["keep_row_focus_on_sort"]
	if ok {
		ct.State.keepRowFocusOnSort = keepRowFocusOnSortIfc.(bool)
	}
	return nil
}

// LoadChartConfig loads chart config from toml config into state struct
func (ct *Cointop) loadChartConfig() error {
	log.Debugf("loadChartConfig()")
	maxChartWidthIfc, ok := ct.config.Chart["max_width"]
	if ok {
		ct.State.maxChartWidth = int(maxChartWidthIfc.(int64))
	}

	chartHeightIfc, ok := ct.config.Chart["height"]
	if ok {
		ct.State.chartHeight = int(chartHeightIfc.(int64))
		ct.State.lastChartHeight = ct.State.chartHeight
	}
	return nil
}

// LoadTableColumnsFromConfig loads preferred coins table columns from config file to struct
func (ct *Cointop) loadTableColumnsFromConfig() error {
	log.Debug("loadTableColumnsFromConfig()")
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
	log.Debug("loadShortcutsFromConfig()")
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
	log.Debug("loadCurrencyFromConfig()")
	if currency, ok := ct.config.Currency.(string); ok {
		ct.State.currencyConversion = strings.ToUpper(currency)
	}
	return nil
}

// LoadDefaultViewFromConfig loads default view from config file to struct
func (ct *Cointop) loadDefaultViewFromConfig() error {
	log.Debug("loadDefaultViewFromConfig()")
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

// LoadDefaultChartRangeFromConfig loads default chart range from config file to struct
func (ct *Cointop) loadDefaultChartRangeFromConfig() error {
	log.Debug("loadDefaultChartRangeFromConfig()")
	if defaultChartRange, ok := ct.config.DefaultChartRange.(string); ok {
		// validate configured value
		_, ok := ct.chartRangesMap[defaultChartRange]
		if !ok {
			return fmt.Errorf("invalid default chart range %q. Valid ranges are: %s", defaultChartRange, strings.Join(ChartRanges(), ","))
		}
		ct.State.defaultChartRange = defaultChartRange
		ct.State.selectedChartRange = defaultChartRange
	}
	return nil
}

// LoadAPIKeysFromConfig loads API keys from config file to struct
func (ct *Cointop) loadAPIKeysFromConfig() error {
	log.Debug("loadAPIKeysFromConfig()")
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
	log.Debug("loadColorschemeFromConfig()")
	if colorscheme, ok := ct.config.Colorscheme.(string); ok {
		ct.colorschemeName = colorscheme
	}

	return nil
}

// LoadRefreshRateFromConfig loads refresh rate from config file to struct
func (ct *Cointop) loadRefreshRateFromConfig() error {
	log.Debug("loadRefreshRateFromConfig()")
	if refreshRate, ok := ct.config.RefreshRate.(int64); ok {
		ct.State.refreshRate = time.Duration(uint(refreshRate)) * time.Second
	}

	return nil
}

// LoadCacheDirFromConfig loads cache dir from config file to struct
func (ct *Cointop) loadCacheDirFromConfig() error {
	log.Debug("loadCacheDirFromConfig()")
	if cacheDir, ok := ct.config.CacheDir.(string); ok {
		ct.State.cacheDir = pathutil.NormalizePath(cacheDir)
	}

	return nil
}

// LoadAPIChoiceFromConfig loads API choices from config file to struct
func (ct *Cointop) loadAPIChoiceFromConfig() error {
	log.Debug("loadAPIKeysFromConfig()")
	apiChoice, ok := ct.config.API.(string)
	if ok {
		apiChoice = strings.TrimSpace(strings.ToLower(apiChoice))
		ct.apiChoice = apiChoice
	}
	return nil
}

// LoadFavoritesFromConfig loads favorites data from config file to struct
func (ct *Cointop) loadFavoritesFromConfig() error {
	log.Debug("loadFavoritesFromConfig()")
	for k, valueIfc := range ct.config.Favorites {
		if k == "character" {
			if favoriteChar, ok := valueIfc.(string); ok {
				if utf8.RuneCountInString(favoriteChar) != 1 {
					return fmt.Errorf("invalid favorite-character. Must be one-character")
				}
				ct.State.favoriteChar = favoriteChar
			}
		}
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
	log.Debug("loadPortfolioFromConfig()")

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
	log.Debug("loadPriceAlertsFromConfig()")
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

// GetColorschemeColors loads colors from colorsheme file to struct
func (ct *Cointop) GetColorschemeColors() (map[string]interface{}, error) {
	log.Debug("GetColorschemeColors()")
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
