package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	bg       = lipgloss.Color("#0D1117")
	ink      = lipgloss.Color("#E6EDF3")
	muted    = lipgloss.Color("#8B949E")
	cyan     = lipgloss.Color("#7DD3FC")
	green    = lipgloss.Color("#A7F3D0")
	yellow   = lipgloss.Color("#FDE68A")
	panel    = lipgloss.Color("#161B22")
	line     = lipgloss.Color("#30363D")
	selected = lipgloss.Color("#1F6FEB")

	headerStyle = lipgloss.NewStyle().
			Foreground(ink).
			Background(bg).
			Bold(true).
			Padding(0, 1)

	warnStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Background(bg).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(muted).
			Background(bg).
			Padding(0, 1)

	filterStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Background(panel).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Foreground(ink).
			Background(panel).
			Border(lipgloss.NormalBorder()).
			BorderForeground(line).
			Padding(1, 2)
)

func tableStyles() table.Styles {
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(line).
		BorderBottom(true).
		Bold(true).
		Foreground(green)
	styles.Selected = styles.Selected.
		Foreground(ink).
		Background(selected).
		Bold(false)
	styles.Cell = styles.Cell.Foreground(ink)
	return styles
}

func columns(width int) []table.Column {
	cmdWidth := width - 74
	remoteWidth := 22
	cmdWidth = width - 98
	if cmdWidth < 24 {
		cmdWidth = 24
	}
	if width < 110 {
		remoteWidth = 16
	}
	containerWidth := 14
	if width < 100 {
		containerWidth = 10
	}
	return []table.Column{
		{Title: "NET", Width: 5},
		{Title: "LOCAL", Width: 22},
		{Title: "REMOTE", Width: remoteWidth},
		{Title: "STATE", Width: 10},
		{Title: "PID", Width: 7},
		{Title: "PROCESS", Width: 16},
		{Title: "CONTAINER", Width: containerWidth},
		{Title: "COMMAND", Width: cmdWidth},
	}
}
