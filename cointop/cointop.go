package cointop

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/miguelmota/cointop/pkg/api"
	"github.com/miguelmota/cointop/pkg/cache"
	"github.com/miguelmota/cointop/pkg/gocui"
	"github.com/miguelmota/cointop/pkg/table"
	"github.com/miguelmota/cointop/pkg/termui"
)

// Cointop cointop
type Cointop struct {
	g                 *gocui.Gui
	marketview        *gocui.View
	chartview         *gocui.View
	chartpoints       [][]termui.Cell
	headersview       *gocui.View
	tableview         *gocui.View
	table             *table.Table
	statusbarview     *gocui.View
	sortdesc          bool
	sortby            string
	api               api.Interface
	allcoins          []*coin
	coins             []*coin
	allcoinsmap       map[string]*coin
	page              int
	perpage           int
	refreshmux        sync.Mutex
	refreshticker     *time.Ticker
	forcerefresh      chan bool
	selectedcoin      *coin
	maxtablewidth     int
	actionsmap        map[string]bool
	shortcutkeys      map[string]string
	config            config // toml config
	searchfield       *gocui.View
	favorites         map[string]bool
	filterByFavorites bool
	savemux           sync.Mutex
	cache             *cache.Cache
	debug             bool
	helpview          *gocui.View
	helpvisible       bool
}

// Run runs cointop
func Run() {
	var debug bool
	if os.Getenv("DEBUG") != "" {
		debug = true
	}
	ct := Cointop{
		api:           api.NewCMC(),
		refreshticker: time.NewTicker(1 * time.Minute),
		sortby:        "rank",
		sortdesc:      false,
		page:          0,
		perpage:       100,
		forcerefresh:  make(chan bool),
		maxtablewidth: 175,
		actionsmap:    actionsMap(),
		shortcutkeys:  defaultShortcuts(),
		favorites:     map[string]bool{},
		cache:         cache.New(1*time.Minute, 2*time.Minute),
		debug:         debug,
	}
	err := ct.setupConfig()
	if err != nil {
		log.Fatal(err)
	}
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
	/*
		ifc, ok, _ := ct.readHardCache(&data, filename)
		if ok {
			// hard cache hit
			if ifc != nil {
				ct.debuglog("hard cache hit")
			}
		}
	*/
	/*
		ifc, ok, _ := ct.readHardCache(&allcoinsmap, "allcoinsmap")
		if ok {
			// hard cache hit
			if ifc != nil {
				ct.debuglog("hard cache hit")
			}
		}
	*/
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
