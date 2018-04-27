package cointop

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	headersview         *gocui.View
	headerviewname      string
	tableview           *gocui.View
	tableviewname       string
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
	favorites           map[string]bool
	filterByFavorites   bool
	savemux             sync.Mutex
	cache               *cache.Cache
	debug               bool
	helpview            *gocui.View
	helpviewname        string
	helpvisible         bool
}

// Run runs cointop
func Run() {
	var ver bool
	flag.BoolVar(&ver, "v", false, "Version")
	flag.Parse()
	if ver {
		fmt.Println("1.0.0")
		return
	}
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	ct := Cointop{
		api:                 api.NewCMC(),
		refreshticker:       time.NewTicker(1 * time.Minute),
		sortby:              "rank",
		sortdesc:            false,
		page:                0,
		perpage:             100,
		forcerefresh:        make(chan bool),
		maxtablewidth:       175,
		actionsmap:          actionsMap(),
		shortcutkeys:        defaultShortcuts(),
		favorites:           map[string]bool{},
		cache:               cache.New(1*time.Minute, 2*time.Minute),
		debug:               debug,
		marketbarviewname:   "market",
		chartviewname:       "chart",
		headerviewname:      "header",
		tableviewname:       "table",
		statusbarviewname:   "statusbar",
		searchfieldviewname: "searchfield",
		helpviewname:        "help",
	}
	err := ct.setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	allcoinsmap := map[string]types.Coin{}
	fcache.Get("allcoinsmap", &allcoinsmap)
	ct.cache.Set("allcoinsmap", allcoinsmap, 10*time.Second)

	var globaldata []float64
	fcache.Get("globaldata", &globaldata)
	ct.cache.Set("globaldata", globaldata, 10*time.Second)

	var market types.GlobalMarketData
	fcache.Get("market", &market)
	ct.cache.Set("market", market, 10*time.Second)

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	ct.g = g
	defer g.Close()
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
	if ct.helpvisible {
		return nil
	}

	return ct.forceQuit()
}

func (ct *Cointop) forceQuit() error {
	return gocui.ErrQuit
}
