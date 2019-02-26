package color

import "github.com/fatih/color"

// Color struct
type Color color.Color

var (
	// Black color
	Black = color.New(color.FgBlack).SprintFunc()
	// BlackBg color
	BlackBg = color.New(color.BgBlack, color.FgWhite).SprintFunc()
	// White color
	White = color.New(color.FgWhite).SprintFunc()
	// WhiteBold bold
	WhiteBold = color.New(color.FgWhite, color.Bold).SprintFunc()
	// Yellow color
	Yellow = color.New(color.FgYellow).SprintFunc()
	// YellowBold color
	YellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	// YellowBg color
	YellowBg = color.New(color.BgYellow, color.FgBlack).SprintFunc()
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
