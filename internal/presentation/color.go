package presentation

import "github.com/fatih/color"

var (
	CreateColor  = color.New(color.FgGreen).SprintFunc()
	UpdateColor  = color.New(color.FgYellow).SprintFunc()
	DeleteColor  = color.New(color.FgRed).SprintFunc()
	ReplaceColor = color.New(color.FgMagenta).SprintFunc()
)
