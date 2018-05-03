package cointop

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miguelmota/cointop/pkg/api"
	types "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/cache"
	"github.com/miguelmota/cointop/pkg/fcache"
	"github.com/miguelmota/cointop/pkg/gocui"
	"github.com/miguelmota/cointop/pkg/table"
	"github.com/miguelmota/cointop/pkg/termui"
)

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
	statusbarview       *gocui.View
	statusbarviewname   string
	sortdesc            bool
	sortby              string
	api                 api.Interface
	allcoins            []*coin
	coins               []*coin
	allcoinsmap         map[string]*coin
	page                int
	perpage             int
	refreshmux          sync.Mutex
	refreshticker       *time.Ticker
	forcerefresh        chan bool
	selectedcoin        *coin
	actionsmap          map[string]bool
	shortcutkeys        map[string]string
	config              config // toml config
	searchfield         *gocui.View
	searchfieldviewname string
	searchfieldvisible  bool
	favorites           map[string]bool
	filterByFavorites   bool
	savemux             sync.Mutex
	cache               *cache.Cache
	debug               bool
	helpview            *gocui.View
	helpviewname        string
	helpvisible         bool
}

// Instance running cointop instance
var Instance *Cointop

// Run runs cointop
func Run() {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	ct := Cointop{
		api:               api.NewCMC(),
		refreshticker:     time.NewTicker(1 * time.Minute),
		sortby:            "rank",
		sortdesc:          false,
		page:              0,
		perpage:           100,
		forcerefresh:      make(chan bool),
		maxtablewidth:     175,
		actionsmap:        actionsMap(),
		shortcutkeys:      defaultShortcuts(),
		favorites:         map[string]bool{},
		cache:             cache.New(1*time.Minute, 2*time.Minute),
		debug:             debug,
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
			"marketcap",
			"24hvolume",
			"1hchange",
			"7dchange",
			"totalsupply",
			"availablesupply",
			"lastupdated",
		},
		statusbarviewname:   "statusbar",
		searchfieldviewname: "searchfield",
		helpviewname:        "help",
	}
	Instance = &ct
	err := ct.setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	allcoinsmap := map[string]types.Coin{}
	coinscachekey := "allcoinsmap"
	fcache.Get(coinscachekey, &allcoinsmap)
	ct.cache.Set(coinscachekey, allcoinsmap, 10*time.Second)

	var globaldata []float64
	chartcachekey := strings.ToLower(fmt.Sprintf("%s_%s", "globaldata", strings.Replace(ct.selectedchartrange, " ", "", -1)))
	fcache.Get(chartcachekey, &globaldata)
	ct.cache.Set(chartcachekey, globaldata, 10*time.Second)

	var market types.GlobalMarketData
	marketcachekey := "market"
	fcache.Get(marketcachekey, &market)
	ct.cache.Set(marketcachekey, market, 10*time.Second)

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	ct.g = g
	defer g.Close()
	g.InputEsc = true
	g.BgColor = gocui.ColorBlack
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

func (ct *Cointop) quit() error {
	if ct.helpvisible || ct.searchfieldvisible {
		return nil
	}

	return ct.forceQuit()
}

func (ct *Cointop) forceQuit() error {
	return gocui.ErrQuit
}

// Exit safely exit application
func Exit() {
	if Instance != nil {
		Instance.g.Close()
	}
}
