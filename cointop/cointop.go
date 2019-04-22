package cointop

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/cointop/common/api"
	"github.com/miguelmota/cointop/cointop/common/api/types"
	"github.com/miguelmota/cointop/cointop/common/filecache"
	"github.com/miguelmota/cointop/cointop/common/table"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// TODO: clean up and optimize codebase

// Cointop cointop
type Cointop struct {
	g                   *gocui.Gui
	marketbarviewname   string
	marketbarview       *gocui.View
	chartview           *gocui.View
	chartviewname       string
	chartpoints         [][]termui.Cell
	chartranges         []string
	chartrangesmap      map[string]time.Duration
	selectedchartrange  string
	headersview         *gocui.View
	headerviewname      string
	tableview           *gocui.View
	tableviewname       string
	tablecolumnorder    []string
	table               *table.Table
	maxtablewidth       int
	portfoliovisible    bool
	visible             bool
	statusbarview       *gocui.View
	statusbarviewname   string
	sortdesc            bool
	sortby              string
	api                 api.Interface
	allcoins            []*coin
	coins               []*coin
	allcoinsslugmap     map[string]*coin
	page                int
	perpage             int
	refreshmux          sync.Mutex
	refreshticker       *time.Ticker
	forcerefresh        chan bool
	selectedcoin        *coin
	actionsmap          map[string]bool
	shortcutkeys        map[string]string
	config              config // toml config
	configFilepath      string
	searchfield         *gocui.View
	searchfieldviewname string
	searchfieldvisible  bool
	// DEPRECATED: favorites by 'symbol' is deprecated because of collisions.
	favoritesbysymbol           map[string]bool
	favorites                   map[string]bool
	filterByFavorites           bool
	savemux                     sync.Mutex
	cache                       *cache.Cache
	debug                       bool
	helpview                    *gocui.View
	helpviewname                string
	helpvisible                 bool
	currencyconversion          string
	convertmenuview             *gocui.View
	convertmenuviewname         string
	convertmenuvisible          bool
	portfolio                   *portfolio
	portfolioupdatemenuview     *gocui.View
	portfolioupdatemenuviewname string
	portfolioupdatemenuvisible  bool
	inputview                   *gocui.View
	inputviewname               string
	defaultView                 string
	apiKeys                     *apiKeys
	limiter                     <-chan time.Time
}

// PortfolioEntry is portfolio entry
type portfolioEntry struct {
	Coin     string
	Holdings float64
}

// Portfolio is portfolio structure
type portfolio struct {
	Entries map[string]*portfolioEntry
}

// Config config options
type Config struct {
	ConfigFilepath string
}

// apiKeys is api keys structure
type apiKeys struct {
	cmc string
}

var defaultConfigPath = "~/.cointop/config"

// NewCointop initializes cointop
func NewCointop(config *Config) *Cointop {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}

	configFilepath := defaultConfigPath
	if config != nil {
		if config.ConfigFilepath != "" {
			configFilepath = config.ConfigFilepath
		}
	}

	ct := &Cointop{
		allcoinsslugmap: make(map[string]*coin),
		allcoins:        []*coin{},
		refreshticker:   time.NewTicker(1 * time.Minute),
		sortby:          "rank",
		page:            0,
		perpage:         100,
		forcerefresh:    make(chan bool),
		maxtablewidth:   175,
		actionsmap:      actionsMap(),
		shortcutkeys:    defaultShortcuts(),
		// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
		favoritesbysymbol: make(map[string]bool),
		favorites:         make(map[string]bool),
		cache:             cache.New(1*time.Minute, 2*time.Minute),
		debug:             debug,
		configFilepath:    configFilepath,
		marketbarviewname: "market",
		chartviewname:     "chart",
		chartranges: []string{
			"1H",
			"6H",
			"24H",
			"3D",
			"7D",
			"1M",
			"3M",
			"6M",
			"1Y",
			"YTD",
			"All Time",
		},
		chartrangesmap: map[string]time.Duration{
			"All Time": time.Duration(24 * 7 * 4 * 12 * 5 * time.Hour),
			"YTD":      time.Duration(1 * time.Second), // this will be calculated
			"1Y":       time.Duration(24 * 7 * 4 * 12 * time.Hour),
			"6M":       time.Duration(24 * 7 * 4 * 6 * time.Hour),
			"3M":       time.Duration(24 * 7 * 4 * 3 * time.Hour),
			"1M":       time.Duration(24 * 7 * 4 * time.Hour),
			"7D":       time.Duration(24 * 7 * time.Hour),
			"3D":       time.Duration(24 * 3 * time.Hour),
			"24H":      time.Duration(24 * time.Hour),
			"6H":       time.Duration(6 * time.Hour),
			"1H":       time.Duration(1 * time.Hour),
		},
		selectedchartrange: "7D",
		headerviewname:     "header",
		tableviewname:      "table",
		tablecolumnorder: []string{
			"rank",
			"name",
			"symbol",
			"price",
			"holdings",
			"balance",
			"marketcap",
			"24hvolume",
			"1hchange",
			"7dchange",
			"totalsupply",
			"availablesupply",
			"percentholdings",
			"lastupdated",
		},
		statusbarviewname:   "statusbar",
		searchfieldviewname: "searchfield",
		helpviewname:        "help",
		convertmenuviewname: "convertmenu",
		currencyconversion:  "USD",
		portfolio: &portfolio{
			Entries: make(map[string]*portfolioEntry, 0),
		},
		portfolioupdatemenuviewname: "portfolioupdatemenu",
		inputviewname:               "input",
		apiKeys:                     new(apiKeys),
		limiter:                     time.Tick(2 * time.Second),
	}

	err := ct.setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	ct.api = api.NewCMC(ct.apiKeys.cmc)

	coinscachekey := "allcoinsslugmap"
	filecache.Get(coinscachekey, &ct.allcoinsslugmap)

	for k := range ct.allcoinsslugmap {
		ct.allcoins = append(ct.allcoins, ct.allcoinsslugmap[k])
	}
	if len(ct.allcoins) > 1 {
		max := len(ct.allcoins)
		if max > 100 {
			max = 100
		}
		ct.sort(ct.sortby, ct.sortdesc, ct.allcoins, false)
		ct.coins = ct.allcoins[0:max]
	}

	// DEPRECATED: favorites by 'symbol' is deprecated because of collisions. Kept for backward compatibility.
	// Here we're doing a lookup based on symbol and setting the favorite to the coin name instead of coin symbol.
	for i := range ct.allcoinsslugmap {
		coin := ct.allcoinsslugmap[i]
		for k := range ct.favoritesbysymbol {
			if coin.Symbol == k {
				ct.favorites[coin.Name] = true
				delete(ct.favoritesbysymbol, k)
			}
		}
	}

	var globaldata []float64
	chartcachekey := strings.ToLower(fmt.Sprintf("%s_%s", "globaldata", strings.Replace(ct.selectedchartrange, " ", "", -1)))
	filecache.Get(chartcachekey, &globaldata)
	ct.cache.Set(chartcachekey, globaldata, 10*time.Second)

	var market types.GlobalMarketData
	marketcachekey := "market"
	filecache.Get(marketcachekey, &market)
	ct.cache.Set(marketcachekey, market, 10*time.Second)
	err = ct.api.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return ct
}

// Run runs cointop
func (ct *Cointop) Run() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	ct.g = g
	defer g.Close()
	g.InputEsc = true
	//g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorWhite
	g.Mouse = true
	g.Highlight = true
	g.SetManagerFunc(ct.layout)
	if err := ct.keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}

// Clean ...
func Clean() {
	tmpPath := "/tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(tmpPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if strings.HasPrefix(f.Name(), "fcache.") {
				file := fmt.Sprintf("%s/%s", tmpPath, f.Name())
				fmt.Printf("removing %s\n", file)
				if err := os.Remove(file); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	fmt.Println("cointop cache has been cleaned")
}

// Reset ...
func Reset() {
	Clean()

	// default config path
	configPath := fmt.Sprintf("%s%s", userHomeDir(), "/.cointop")
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		fmt.Printf("removing %s\n", configPath)
		if err := os.RemoveAll(configPath); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("cointop has been reset")
}
