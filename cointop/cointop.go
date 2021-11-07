package cointop

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cointop-sh/cointop/pkg/api"
	"github.com/cointop-sh/cointop/pkg/api/types"
	"github.com/cointop-sh/cointop/pkg/cache"
	"github.com/cointop-sh/cointop/pkg/filecache"
	"github.com/cointop-sh/cointop/pkg/gocui"
	"github.com/cointop-sh/cointop/pkg/pathutil"
	"github.com/cointop-sh/cointop/pkg/table"
	"github.com/cointop-sh/cointop/pkg/ui"

	log "github.com/sirupsen/logrus"
)

// TODO: clean up and optimize codebase

// Views are all views in cointop
type Views struct {
	Chart       *ChartView
	Table       *TableView
	TableHeader *TableHeaderView
	Marketbar   *MarketbarView
	SearchField *SearchFieldView
	Statusbar   *StatusbarView
	Menu        *MenuView
	Input       *InputView
}

// State is the state preferences of cointop
type State struct {
	allCoins           []*Coin
	allCoinsSlugMap    sync.Map
	cacheDir           string
	coins              []*Coin
	chartPoints        [][]rune
	currencyConversion string
	coinsTableColumns  []string
	convertMenuVisible bool
	defaultView        string
	defaultChartRange  string
	maxChartWidth      int
	columnLookup       []string

	favorites                  map[string]bool
	favoritesTableColumns      []string
	favoriteChar               string
	helpVisible                bool
	hideMarketbar              bool
	hideChart                  bool
	hideTable                  bool
	hideStatusbar              bool
	hidePortfolioBalances      bool
	keepRowFocusOnSort         bool
	lastSelectedRowIndex       int
	marketBarHeight            int
	maxPages                   int
	page                       int
	perPage                    int
	portfolio                  *Portfolio
	portfolioUpdateMenuVisible bool
	portfolioTableColumns      []string
	refreshRate                time.Duration
	running                    bool
	searchFieldVisible         bool
	lastSearchQuery            string
	selectedCoin               *Coin
	selectedChartRange         string
	selectedView               string
	lastSelectedView           string
	shortcutKeys               map[string]string
	sortDesc                   bool
	sortBy                     string
	tableOffsetX               int
	onlyTable                  bool
	onlyChart                  bool
	tableColumnWidths          sync.Map
	tableColumnAlignLeft       sync.Map
	chartHeight                int
	lastChartHeight            int
	priceAlerts                *PriceAlerts
	priceAlertEditID           string
	priceAlertNewID            string

	compactNotation          bool
	tableCompactNotation     bool
	favoritesCompactNotation bool
	portfolioCompactNotation bool
	enableMouse              bool
}

// Cointop cointop
type Cointop struct {
	g               *gocui.Gui
	ui              *ui.UI
	ActionsMap      map[string]bool
	apiKeys         *APIKeys
	cache           *cache.Cache
	colorsDir       string
	config          ConfigFileConfig
	configFilepath  string
	api             api.Interface
	apiChoice       string
	chartRanges     []string
	chartRangesMap  map[string]time.Duration
	colorschemeName string
	colorscheme     *Colorscheme
	filecache       *filecache.FileCache
	logfile         *os.File
	forceRefresh    chan bool
	limiter         <-chan time.Time
	maxTableWidth   int
	refreshMux      sync.Mutex
	refreshTicker   *time.Ticker
	saveMux         sync.Mutex
	State           *State
	table           *table.Table
	Views           *Views
}

// PortfolioEntry is portfolio entry
type PortfolioEntry struct {
	Coin        string
	Holdings    float64
	BuyPrice    float64
	BuyCurrency string
}

// Portfolio is portfolio structure
type Portfolio struct {
	Entries map[string]*PortfolioEntry
}

// PriceAlert is price alert structure
type PriceAlert struct {
	ID          string
	CoinName    string
	TargetPrice float64
	Operator    string
	Frequency   string
	CreatedAt   string
	Expired     bool
}

// PriceAlerts is price alerts structure
type PriceAlerts struct {
	Entries      []*PriceAlert
	SoundEnabled bool
}

// Config config options
type Config struct {
	APIChoice             string
	CacheDir              string
	ColorsDir             string
	Colorscheme           string
	ConfigFilepath        string
	CoinMarketCapAPIKey   string
	NoPrompts             bool
	HideMarketbar         bool
	HideChart             bool
	HideTable             bool
	HideStatusbar         bool
	HidePortfolioBalances bool
	NoCache               bool
	OnlyTable             bool
	OnlyChart             bool
	RefreshRate           *uint
	PerPage               uint
	MaxPages              uint
}

// APIKeys is api keys structure
type APIKeys struct {
	cmc string
}

// DefaultCurrency ...
var DefaultCurrency = "USD"

// DefaultChartRange ...
var DefaultChartRange = "1Y"

// DefaultCompactNotation ...
var DefaultCompactNotation = false

// DefaultEnableMouse ...
var DefaultEnableMouse = true

// DefaultMaxChartWidth ...
var DefaultMaxChartWidth = 175

// DefaultChartHeight ...
var DefaultChartHeight = 10

// DefaultSortBy ...
var DefaultSortBy = "rank"

// DefaultPerPage ...
var DefaultPerPage = uint(100)

// DefaultMaxPages ...
var DefaultMaxPages = uint(10)

// DefaultColorscheme ...
var DefaultColorscheme = "cointop"

// DefaultConfigFilepath ...
var DefaultConfigFilepath = pathutil.NormalizePath(":PREFERRED_CONFIG_HOME:/cointop/config.toml")

// DefaultCacheDir ...
var DefaultCacheDir = filecache.DefaultCacheDir

// DefaultColorsDir ...
var DefaultColorsDir = fmt.Sprintf("%s/colors", DefaultConfigFilepath)

// DefaultFavoriteChar ...
var DefaultFavoriteChar = "*"

// NewCointop initializes cointop
func NewCointop(config *Config) (*Cointop, error) {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}

	if config == nil {
		config = &Config{}
	}

	configFilepath := DefaultConfigFilepath
	if config.ConfigFilepath != "" {
		configFilepath = config.ConfigFilepath
	}

	perPage := DefaultPerPage
	if config.PerPage > 0 {
		perPage = config.PerPage
	}

	maxPages := DefaultMaxPages
	if config.MaxPages > 0 {
		maxPages = config.MaxPages
	}

	ct := &Cointop{
		// defaults
		apiChoice:      CoinGecko,
		apiKeys:        new(APIKeys),
		forceRefresh:   make(chan bool),
		maxTableWidth:  175,
		ActionsMap:     ActionsMap(),
		cache:          cache.New(1*time.Minute, 2*time.Minute),
		colorsDir:      config.ColorsDir,
		configFilepath: configFilepath,
		chartRanges:    ChartRanges(),
		chartRangesMap: ChartRangesMap(),
		limiter:        time.NewTicker(2 * time.Second).C,
		filecache:      nil,
		State: &State{
			allCoins:              []*Coin{},
			cacheDir:              DefaultCacheDir,
			coinsTableColumns:     DefaultCoinTableHeaders,
			currencyConversion:    DefaultCurrency,
			defaultChartRange:     DefaultChartRange,
			maxChartWidth:         DefaultMaxChartWidth,
			favorites:             make(map[string]bool),
			favoritesTableColumns: DefaultCoinTableHeaders,
			favoriteChar:          DefaultFavoriteChar,
			hideMarketbar:         config.HideMarketbar,
			hideChart:             config.HideChart,
			hideTable:             config.HideTable,
			hideStatusbar:         config.HideStatusbar,
			hidePortfolioBalances: config.HidePortfolioBalances,
			keepRowFocusOnSort:    false,
			marketBarHeight:       1,
			maxPages:              int(maxPages),
			onlyTable:             config.OnlyTable,
			onlyChart:             config.OnlyChart,
			refreshRate:           60 * time.Second,
			selectedChartRange:    DefaultChartRange,
			shortcutKeys:          DefaultShortcuts(),
			sortBy:                DefaultSortBy,
			page:                  0,
			perPage:               int(perPage),
			portfolio: &Portfolio{
				Entries: make(map[string]*PortfolioEntry),
			},
			portfolioTableColumns: DefaultPortfolioTableHeaders,
			chartHeight:           DefaultChartHeight,
			lastChartHeight:       DefaultChartHeight,
			tableOffsetX:          0,
			tableColumnWidths:     sync.Map{},
			tableColumnAlignLeft:  sync.Map{},
			priceAlerts: &PriceAlerts{
				Entries:      make([]*PriceAlert, 0),
				SoundEnabled: true,
			},
			compactNotation:          DefaultCompactNotation,
			enableMouse:              DefaultEnableMouse,
			tableCompactNotation:     DefaultCompactNotation,
			favoritesCompactNotation: DefaultCompactNotation,
			portfolioCompactNotation: DefaultCompactNotation,
		},
		Views: &Views{
			Chart:       NewChartView(),
			Table:       NewTableView(),
			TableHeader: NewTableHeaderView(),
			Marketbar:   NewMarketbarView(),
			SearchField: NewSearchFieldView(),
			Statusbar:   NewStatusbarView(),
			Menu:        NewMenuView(),
			Input:       NewInputView(),
		},
	}
	ct.initlog()

	err := ct.SetupConfig()
	if err != nil {
		return nil, err
	}

	ct.cache.Set("onlyTable", ct.State.onlyTable, cache.NoExpiration)
	if ct.State.onlyTable && ct.State.onlyChart {
		ct.State.onlyChart = false
	}
	ct.cache.Set("onlyChart", ct.State.onlyChart, cache.NoExpiration)
	ct.cache.Set("hideMarketbar", ct.State.hideMarketbar, cache.NoExpiration)
	ct.cache.Set("hideChart", ct.State.hideChart, cache.NoExpiration)
	ct.cache.Set("hideTable", ct.State.hideTable, cache.NoExpiration)
	ct.cache.Set("hideStatusbar", ct.State.hideStatusbar, cache.NoExpiration)

	if config.RefreshRate != nil {
		ct.State.refreshRate = time.Duration(*config.RefreshRate) * time.Second
	}

	if ct.State.refreshRate == 0 {
		ct.refreshTicker = time.NewTicker(time.Duration(1))
		ct.refreshTicker.Stop()
	} else {
		ct.refreshTicker = time.NewTicker(ct.State.refreshRate)
	}

	if config.CacheDir != "" {
		ct.State.cacheDir = pathutil.NormalizePath(config.CacheDir)
		if err := ct.SaveConfig(); err != nil {
			return nil, err
		}
	}

	if !config.NoCache {
		// each custom config file has it's own file cache
		hash := sha256.Sum256([]byte(ct.ConfigFilePath()))
		fcache, err := filecache.NewFileCache(&filecache.Config{
			CacheDir: ct.State.cacheDir,
			Prefix:   fmt.Sprintf("%x", hash[0:4]),
		})
		if err != nil {
			fmt.Printf("error: %s\nyou may change the cache directory with --cache-dir flag.\nproceeding without filecache.\n", err)
		}

		ct.filecache = fcache
	}

	// prompt for CoinMarketCap api key if not found
	if config.CoinMarketCapAPIKey != "" {
		ct.apiKeys.cmc = config.CoinMarketCapAPIKey
		if err := ct.SaveConfig(); err != nil {
			return nil, err
		}
	}

	if config.Colorscheme != "" {
		ct.colorschemeName = config.Colorscheme
	}

	colors, err := ct.GetColorschemeColors()
	if err != nil {
		return nil, err
	}
	ct.colorscheme = NewColorscheme(colors)

	if config.APIChoice != "" {
		ct.apiChoice = config.APIChoice
		if err := ct.SaveConfig(); err != nil {
			return nil, err
		}
	}

	if ct.apiChoice == CoinMarketCap && ct.apiKeys.cmc == "" {
		apiKey := os.Getenv("CMC_PRO_API_KEY")
		if apiKey == "" {
			if !config.NoPrompts {
				apiKey, err = ct.ReadAPIKeyFromStdin("CoinMarketCap Pro")
				if err != nil {
					return nil, err
				}

				ct.apiKeys.cmc = apiKey
			}
		} else {
			ct.apiKeys.cmc = apiKey
		}

		if err := ct.SaveConfig(); err != nil {
			return nil, err
		}
	}

	if ct.apiChoice == CoinMarketCap {
		ct.api = api.NewCMC(ct.apiKeys.cmc)
	} else if ct.apiChoice == CoinGecko {
		ct.api = api.NewCG(perPage, maxPages)
	} else {
		return nil, ErrInvalidAPIChoice
	}

	allCoinsSlugMap := make(map[string]*Coin)
	coinscachekey := ct.CacheKey("allCoinsSlugMap")
	if ct.filecache != nil {
		ct.filecache.Get(coinscachekey, &allCoinsSlugMap)
	}

	// fix for https://github.com/cointop-sh/cointop/issues/59
	// can remove this after everyone has cleared their cache
	for _, v := range allCoinsSlugMap {
		// Some APIs returns rank 0 for new coins
		// or coins with low market cap data so we need to put them
		// at the end of the list.
		if v.Rank == 0 {
			v.Rank = 10000
		}
	}

	for k, v := range allCoinsSlugMap {
		ct.State.allCoinsSlugMap.Store(k, v)
	}

	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		if coin, ok := value.(*Coin); ok {
			ct.State.allCoins = append(ct.State.allCoins, coin)
		}
		return true
	})

	if len(ct.State.allCoins) > 1 {
		max := len(ct.State.allCoins)
		if max > 100 {
			max = 100
		}
		ct.Sort(ct.State.sortBy, ct.State.sortDesc, ct.State.allCoins, false)
		ct.State.coins = ct.State.allCoins[0:max]
	}

	var globaldata []float64
	chartcachekey := ct.CompositeCacheKey("globaldata", "", "", ct.State.selectedChartRange)
	if ct.filecache != nil {
		ct.filecache.Get(chartcachekey, &globaldata)
	}
	ct.cache.Set(chartcachekey, globaldata, 10*time.Second)

	var market types.GlobalMarketData
	marketcachekey := ct.CacheKey("market")
	if ct.filecache != nil {
		ct.filecache.Get(marketcachekey, &market)
	}
	ct.cache.Set(marketcachekey, market, 10*time.Second)

	// TODO: notify offline status in status bar
	/*
		if err := ct.api.Ping(); err != nil {
			return nil, err
		}
	*/
	return ct, nil
}

// Run runs cointop
func (ct *Cointop) Run() error {
	log.Debug("Run()")
	ui, err := ui.NewUI()
	if err != nil {
		return err
	}

	ui.SetStyle(ct.colorscheme.BaseStyle())
	ct.ui = ui
	ct.g = ui.GetGocui()
	defer ui.Close()

	ui.SetInputEsc(true)
	ui.SetMouse(ct.State.enableMouse)
	ui.SetHighlight(true)
	ui.SetManagerFunc(ct.layout)
	if err := ct.SetKeybindings(); err != nil {
		return fmt.Errorf("keybindings: %v", err)
	}

	go ct.PriceAlertWatcher()
	ct.State.running = true
	if err := ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("main loop: %v", err)
	}

	return nil
}

// IsRunning returns true if cointop is running
func (ct *Cointop) IsRunning() bool {
	return ct.State.running
}

// CleanConfig is the config for the clean function
type CleanConfig struct {
	Log      bool
	CacheDir string
}

// Clean removes cache files
func (ct *Cointop) Clean(config *CleanConfig) error {
	if config == nil {
		config = &CleanConfig{}
	}
	cacheDir := DefaultCacheDir
	if config.CacheDir != "" {
		cacheDir = pathutil.NormalizePath(config.CacheDir)
	} else if ct.State.cacheDir != "" {
		cacheDir = ct.State.cacheDir
	}

	cacheCleaned := false

	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(cacheDir)
		if err != nil {
			return err
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), "fcache.") {
				file := fmt.Sprintf("%s/%s", cacheDir, f.Name())
				if config.Log {
					fmt.Printf("removing %s\n", file)
				}
				if err := os.Remove(file); err != nil {
					return err
				}

				cacheCleaned = true
			}
		}
	}

	if config.Log {
		if cacheCleaned {
			fmt.Println("cointop cache has been cleaned")
		}
	}

	return nil
}

// ResetConfig is the config for the reset function
type ResetConfig struct {
	Log      bool
	CacheDir string
}

// Reset removes configuration and cache files
func (ct *Cointop) Reset(config *ResetConfig) error {
	if config == nil {
		config = &ResetConfig{}
	}

	if err := ct.Clean(&CleanConfig{
		CacheDir: config.CacheDir,
		Log:      config.Log,
	}); err != nil {
		return err
	}

	configDeleted := false

	for _, configPath := range PossibleConfigPaths {
		normalizedPath := pathutil.NormalizePath(configPath)
		if _, err := os.Stat(normalizedPath); !os.IsNotExist(err) {
			if config.Log {
				fmt.Printf("removing %s\n", normalizedPath)
			}
			if err := os.RemoveAll(normalizedPath); err != nil {
				return err
			}

			configDeleted = true
		}
	}

	if config.Log {
		if configDeleted {
			fmt.Println("cointop has been reset")
		}
	}

	return nil
}

// ColorschemeHelpString ...
func ColorschemeHelpString() string {
	return fmt.Sprintf("To install standard themes, do:\n\ngit clone git@github.com:cointop-sh/colors.git %s\n\nSee git.io/cointop#colorschemes for more info.", pathutil.NormalizePath(":PREFERRED_CONFIG_HOME:/cointop/colors"))
}
