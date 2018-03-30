package color

import "github.com/fatih/color"

var (
	White = color.New(color.FgWhite).SprintFunc()
	Green = color.New(color.FgGreen).SprintFunc()
	Red   = color.New(color.FgRed).SprintFunc()
	Cyan  = color.New(color.FgCyan).SprintFunc()
)
