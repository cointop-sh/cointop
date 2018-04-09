package cointop

import (
	"log"
	"sync"
	"time"

	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/api"
	apt "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/table"
)

// Cointop cointop
type Cointop struct {
	g             *gocui.Gui
	marketview    *gocui.View
	chartview     *gocui.View
	chartpoints   [][]termui.Cell
	headersview   *gocui.View
	tableview     *gocui.View
	table         *table.Table
	statusbarview *gocui.View
	sortdesc      bool
	sortby        string
	api           api.Interface
	allcoins      []*apt.Coin
	coins         []*apt.Coin
	allcoinsmap   map[string]apt.Coin
	page          int
	perpage       int
	refreshmux    sync.Mutex
	refreshticker *time.Ticker
	forcerefresh  chan bool
	selectedcoin  *apt.Coin
	maxtablewidth int
	shortcutkeys  map[string]string
	config        config // toml config
	searchfield   *gocui.View
}

// Run runs cointop
func Run() {
	ct := Cointop{
		api:           api.NewCMC(),
		refreshticker: time.NewTicker(1 * time.Minute),
		sortby:        "rank",
		sortdesc:      false,
		page:          0,
		perpage:       100,
		forcerefresh:  make(chan bool),
		maxtablewidth: 175,
		shortcutkeys:  defaultShortcuts(),
	}
	_ = ct.setupConfig()
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	ct.g = g
	defer g.Close()
	g.Cursor = true
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
	return gocui.ErrQuit
}
