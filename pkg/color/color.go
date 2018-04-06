package color

import "github.com/fatih/color"

var (
	// Black color
	Black = color.New(color.FgBlack).SprintFunc()
	// White color
	White = color.New(color.FgWhite).SprintFunc()
	// WhiteBold bold
	WhiteBold = color.New(color.FgWhite, color.Bold).SprintFunc()
	// Green color
	Green = color.New(color.FgGreen).SprintFunc()
	// GreenBg color
	GreenBg = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	// Red color
	Red = color.New(color.FgRed).SprintFunc()
	// Cyan color
	Cyan = color.New(color.FgCyan).SprintFunc()
	// CyanBg color
	CyanBg = color.New(color.BgCyan, color.FgBlack).SprintFunc()
	// Blue color
	Blue = color.New(color.FgBlue).SprintFunc()
	// BlueBg color
	BlueBg = color.New(color.BgBlue).SprintFunc()
)
