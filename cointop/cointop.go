package cointop

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miguelmota/cointop/cointop/common/api"
	"github.com/miguelmota/cointop/cointop/common/api/types"
	"github.com/miguelmota/cointop/cointop/common/filecache"
	"github.com/miguelmota/cointop/cointop/common/gizak/termui"
	"github.com/miguelmota/cointop/cointop/common/pathutil"
	"github.com/miguelmota/cointop/cointop/common/table"
	"github.com/miguelmota/gocui"
	"github.com/patrickmn/go-cache"
)

// TODO: clean up and optimize codebase

// ErrInvalidAPIChoice is error for invalid API choice
var ErrInvalidAPIChoice = errors.New("Invalid API choice")

// Views are all views in cointop
type Views struct {
	Chart               *ChartView
	Table               *TableView
	TableHeader         *TableHeaderView
	Marketbar           *MarketbarView
	SearchField         *SearchFieldView
	Statusbar           *StatusbarView
	Help                *HelpView
	ConvertMenu         *ConvertMenuView
	Input               *InputView
	PortfolioUpdateMenu *PortfolioUpdateMenuView
}

// State is the state preferences of cointop
type State struct {
	allCoins           []*Coin
	allCoinsSlugMap    sync.Map
	coins              []*Coin
	chartPoints        [][]termui.Cell
	currencyConversion string
	convertMenuVisible bool
	defaultView        string

	// DEPRECATED: favorites by 'symbol' is deprecated because of collisions.
	favoritesBySymbol map[string]bool

	favorites                  map[string]bool
	filterByFavorites          bool
	helpVisible                bool
	hideMarketbar              bool
	hideChart                  bool
	hideStatusbar              bool
	lastSelectedRowIndex       int
	page                       int
	perPage                    int
	portfolio                  *Portfolio
	portfolioVisible           bool
	portfolioUpdateMenuVisible bool
	refreshRate                time.Duration
	running                    bool
	searchFieldVisible         bool
	selectedCoin               *Coin
	selectedChartRange         string
	shortcutKeys               map[string]string
	sortDesc                   bool
	sortBy                     string
	onlyTable                  bool
	chartHeight                int
}

// Cointop cointop
type Cointop struct {
	g                *gocui.Gui
	ActionsMap       map[string]bool
	apiKeys          *APIKeys
	cache            *cache.Cache
	config           config // toml config
	configFilepath   string
	api              api.Interface
	apiChoice        string
	chartRanges      []string
	chartRangesMap   map[string]time.Duration
	colorschemeName  string
	colorscheme      *Colorscheme
	debug            bool
	filecache        *filecache.FileCache
	forceRefresh     chan bool
	limiter          <-chan time.Time
	maxTableWidth    int
	refreshMux       sync.Mutex
	refreshTicker    *time.Ticker
	saveMux          sync.Mutex
	State            *State
	table            *table.Table
	TableColumnOrder []string
	Views            *Views
}

// CoinMarketCap is API choice
var CoinMarketCap = "coinmarketcap"

// CoinGecko is API choice
var CoinGecko = "coingecko"

// PortfolioEntry is portfolio entry
type PortfolioEntry struct {
	Coin     string
	Holdings float64
}

// Portfolio is portfolio structure
type Portfolio struct {
	Entries map[string]*PortfolioEntry
}

// Config config options
type Config struct {
	APIChoice           string
	CacheDir            string
	Colorscheme         string
	ConfigFilepath      string
	CoinMarketCapAPIKey string
	NoPrompts           bool
	HideMarketbar       bool
	HideChart           bool
	HideStatusbar       bool
	NoCache             bool
	OnlyTable           bool
	RefreshRate         *uint
	PerPage             uint
}

// APIKeys is api keys structure
type APIKeys struct {
	cmc string
}

// DefaultPerPage ...
var DefaultPerPage uint = 100

// DefaultColorscheme ...
var DefaultColorscheme = "cointop"

// DefaultConfigFilepath ...
var DefaultConfigFilepath = pathutil.NormalizePath(":PREFERRED_CONFIG_HOME:/cointop/config.toml")

// DefaultCacheDir ...
var DefaultCacheDir = filecache.DefaultCacheDir

// NewCointop initializes cointop
func NewCointop(config *Config) (*Cointop, error) {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}

	if config == nil {
		config = &Config{}
	}

	configFilepath := DefaultConfigFilepath
	if config.ConfigFilepath != "" {
		configFilepath = config.ConfigFilepath
	}

	var fcache *filecache.FileCache
	if !config.NoCache {
		fcache = filecache.NewFileCache(&filecache.Config{
			CacheDir: config.CacheDir,
		})
	}

	perPage := DefaultPerPage
	if config.PerPage != 0 {
		perPage = config.PerPage
	}

	ct := &Cointop{
		apiChoice:      CoinGecko,
		apiKeys:        new(APIKeys),
		forceRefresh:   make(chan bool),
		maxTableWidth:  175,
		ActionsMap:     ActionsMap(),
		cache:          cache.New(1*time.Minute, 2*time.Minute),
		configFilepath: configFilepath,
		chartRanges:    ChartRanges(),
		debug:          debug,
		chartRangesMap: ChartRangesMap(),
		limiter:        time.Tick(2 * time.Second),
		filecache:      fcache,
		State: &State{
			allCoins:           []*Coin{},
			currencyConversion: "USD",
			// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
			favoritesBySymbol:  make(map[string]bool),
			favorites:          make(map[string]bool),
			hideMarketbar:      config.HideMarketbar,
			hideChart:          config.HideChart,
			hideStatusbar:      config.HideStatusbar,
			onlyTable:          config.OnlyTable,
			refreshRate:        60 * time.Second,
			selectedChartRange: "7D",
			shortcutKeys:       DefaultShortcuts(),
			sortBy:             "rank",
			page:               0,
			perPage:            int(perPage),
			portfolio: &Portfolio{
				Entries: make(map[string]*PortfolioEntry, 0),
			},
			chartHeight: 10,
		},
		TableColumnOrder: TableColumnOrder(),
		Views: &Views{
			Chart:               NewChartView(),
			Table:               NewTableView(),
			TableHeader:         NewTableHeaderView(),
			Marketbar:           NewMarketbarView(),
			SearchField:         NewSearchFieldView(),
			Statusbar:           NewStatusbarView(),
			Help:                NewHelpView(),
			ConvertMenu:         NewConvertMenuView(),
			Input:               NewInputView(),
			PortfolioUpdateMenu: NewPortfolioUpdateMenuView(),
		},
	}

	err := ct.SetupConfig()
	if err != nil {
		return nil, err
	}

	ct.cache.Set("onlyTable", ct.State.onlyTable, cache.NoExpiration)
	ct.cache.Set("hideMarketbar", ct.State.hideMarketbar, cache.NoExpiration)
	ct.cache.Set("hideChart", ct.State.hideChart, cache.NoExpiration)
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

	colors, err := ct.getColorschemeColors()
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

	if ct.apiChoice == CoinGecko {
		ct.State.selectedChartRange = "1Y"
	}

	if ct.apiChoice == CoinMarketCap {
		ct.api = api.NewCMC(ct.apiKeys.cmc)
	} else if ct.apiChoice == CoinGecko {
		ct.api = api.NewCG()
	} else {
		return nil, ErrInvalidAPIChoice
	}

	allCoinsSlugMap := make(map[string]*Coin)
	coinscachekey := ct.CacheKey("allCoinsSlugMap")
	if ct.filecache != nil {
		ct.filecache.Get(coinscachekey, &allCoinsSlugMap)
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

	// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
	// Here we're doing a lookup based on symbol and setting the favorite to the coin name instead of coin symbol.
	ct.State.allCoinsSlugMap.Range(func(key, value interface{}) bool {
		if coin, ok := value.(*Coin); ok {
			for k := range ct.State.favoritesBySymbol {
				if coin.Symbol == k {
					ct.State.favorites[coin.Name] = true
					delete(ct.State.favoritesBySymbol, k)
				}
			}
		}

		return true
	})

	var globaldata []float64
	chartcachekey := ct.CacheKey(fmt.Sprintf("%s_%s", "globaldata", strings.Replace(ct.State.selectedChartRange, " ", "", -1)))
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
	ct.debuglog("run()")
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return fmt.Errorf("new gocui: %v", err)
	}

	g.FgColor = ct.colorscheme.BaseFg()
	g.BgColor = ct.colorscheme.BaseBg()
	ct.g = g
	defer g.Close()

	g.InputEsc = true
	g.Mouse = true
	g.Highlight = true
	g.SetManagerFunc(ct.layout)
	if err := ct.Keybindings(g); err != nil {
		return fmt.Errorf("keybindings: %v", err)
	}

	ct.State.running = true
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
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
func Clean(config *CleanConfig) error {
	if config == nil {
		config = &CleanConfig{}
	}

	cacheCleaned := false

	cacheDir := DefaultCacheDir
	if config.CacheDir != "" {
		cacheDir = config.CacheDir
	}

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
func Reset(config *ResetConfig) error {
	if config == nil {
		config = &ResetConfig{}
	}

	if err := Clean(&CleanConfig{
		CacheDir: config.CacheDir,
		Log:      config.Log,
	}); err != nil {
		return err
	}

	configDeleted := false

	for _, configPath := range possibleConfigPaths {
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
