package cointop

import (
	"log"
	"sync"
	"time"

	"github.com/gizak/termui"
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/api"
	apitypes "github.com/miguelmota/cointop/pkg/api/types"
	"github.com/miguelmota/cointop/pkg/table"
)

var (
	oneMinute int64 = 60
	oneHour         = oneMinute * 60
	oneDay          = oneHour * 24
	oneWeek         = oneDay * 7
	oneMonth        = oneDay * 30
	oneYear         = oneDay * 365
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
	statusview    *gocui.View
	sortdesc      bool
	sortby        string
	api           api.Interface
	allcoins      []*apitypes.Coin
	coins         []*apitypes.Coin
	allcoinsmap   map[string]apitypes.Coin
	page          int
	perpage       int
	refreshmux    sync.Mutex
	refreshticker *time.Ticker
}

// Run runs cointop
func Run() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true
	g.Highlight = true
	ct := Cointop{
		g:             g,
		api:           api.NewCMC(),
		refreshticker: time.NewTicker(1 * time.Minute),
		sortby:        "rank",
		sortdesc:      false,
		page:          0,
		perpage:       100,
	}
	g.SetManagerFunc(ct.layout)
	if err := ct.keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}
