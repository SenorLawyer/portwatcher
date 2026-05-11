package tui

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/SenorLawyer/portwatcher/internal/app"
	"github.com/SenorLawyer/portwatcher/internal/config"
	"github.com/SenorLawyer/portwatcher/internal/history"
	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

type actor interface {
	Kill(pid int32, force bool) error
	Copy(value string) error
	Open(url string) error
}

type Model struct {
	service  *app.Service
	actions  actor
	cfg      config.Config
	table    table.Model
	spinner  spinner.Model
	snapshot scanner.Snapshot
	events   []history.Event
	filter   string
	status   string
	err      error
	width    int
	height   int
	mode     mode
	confirm  confirmState
}

type mode int

const (
	modeTable mode = iota
	modeFilter
	modeHistory
	modeConfirm
)

type confirmState struct {
	force bool
	entry scanner.PortEntry
}

type refreshMsg struct {
	snapshot scanner.Snapshot
	events   []history.Event
	err      error
}

type tickMsg time.Time

func New(service *app.Service, actions actor, cfg config.Config) Model {
	sp := spinner.New(spinner.WithSpinner(spinner.MiniDot))
	t := table.New(
		table.WithFocused(true),
		table.WithColumns(columns(120)),
		table.WithRows(nil),
		table.WithHeight(10),
	)
	t.SetStyles(tableStyles())
	return Model{
		service: service,
		actions: actions,
		cfg:     cfg,
		table:   t,
		spinner: sp,
		status:  "starting",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.refresh(), tick(m.cfg.Interval))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetColumns(columns(msg.Width))
		m.table.SetHeight(max(4, msg.Height-8))
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tickMsg:
		return m, tea.Batch(m.refresh(), tick(m.cfg.Interval))
	case refreshMsg:
		m.err = msg.err
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.snapshot = msg.snapshot
		m.events = append(msg.events, m.events...)
		m.status = fmt.Sprintf("%d ports scanned in %s", len(msg.snapshot.Ports), msg.snapshot.Duration)
		m.rebuildRows()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "portwatch"
	}
	header := headerStyle.Width(m.width).Render(fmt.Sprintf("PORTWATCH  %s  %s", m.spinner.View(), m.status))
	if len(m.snapshot.Warnings) > 0 {
		header = lipgloss.JoinVertical(lipgloss.Left, header, warnStyle.Width(m.width).Render(strings.Join(m.snapshot.Warnings, " | ")))
	}
	body := m.table.View()
	if m.mode == modeHistory {
		body = m.historyView()
	}
	if m.mode == modeConfirm {
		body = m.confirmView()
	}
	if m.mode == modeFilter {
		body = filterStyle.Width(m.width).Render("/"+m.filter) + "\n" + body
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footerStyle.Width(m.width).Render(m.footer()))
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.mode == modeFilter {
		switch msg.Type {
		case tea.KeyEsc, tea.KeyEnter:
			m.mode = modeTable
		case tea.KeyBackspace:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.rebuildRows()
			}
		default:
			if msg.String() != "" && len(msg.String()) == 1 {
				m.filter += msg.String()
				m.rebuildRows()
			}
		}
		return m, nil
	}
	if m.mode == modeConfirm {
		switch msg.String() {
		case "y", "Y":
			err := m.actions.Kill(m.confirm.entry.PID, m.confirm.force)
			m.mode = modeTable
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.status = "signal sent"
			return m, m.refresh()
		default:
			m.mode = modeTable
			return m, nil
		}
	}

	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
		return m, tea.Quit
	case msg.String() == "r":
		return m, m.refresh()
	case msg.String() == "/":
		m.mode = modeFilter
		return m, nil
	case msg.String() == "esc":
		m.filter = ""
		m.rebuildRows()
		m.mode = modeTable
		return m, nil
	case msg.String() == "h":
		if m.mode == modeHistory {
			m.mode = modeTable
		} else {
			m.mode = modeHistory
		}
		return m, nil
	case msg.String() == "k" || msg.String() == "K":
		entry, ok := m.selected()
		if !ok {
			return m, nil
		}
		m.confirm = confirmState{force: msg.String() == "K", entry: entry}
		m.mode = modeConfirm
		return m, nil
	case msg.String() == "c":
		if entry, ok := m.selected(); ok {
			m.status = report(m.actions.Copy(rowSummary(entry)), "copied row")
		}
		return m, nil
	case msg.String() == "p":
		if entry, ok := m.selected(); ok {
			m.status = report(m.actions.Copy(strconv.Itoa(int(entry.Port))), "copied port")
		}
		return m, nil
	case msg.String() == "o":
		if entry, ok := m.selected(); ok && entry.Protocol == "tcp" {
			target := (&url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", entry.Port)}).String()
			m.status = report(m.actions.Open(target), "opened "+target)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) refresh() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), maxDuration(2*m.cfg.Interval, 3*time.Second))
		defer cancel()
		snap, events, err := m.service.Refresh(ctx)
		return refreshMsg{snapshot: snap, events: events, err: err}
	}
}

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m *Model) rebuildRows() {
	rows := []table.Row{}
	for _, p := range m.snapshot.Ports {
		if !matches(p, m.filter) {
			continue
		}
		rows = append(rows, table.Row{
			p.Protocol,
			p.Address,
			p.Remote,
			p.State,
			strconv.Itoa(int(p.PID)),
			p.Process,
			p.Container,
			p.Command,
		})
	}
	m.table.SetRows(rows)
}

func (m Model) selected() (scanner.PortEntry, bool) {
	row := m.table.SelectedRow()
	if len(row) == 0 {
		return scanner.PortEntry{}, false
	}
	for _, p := range m.snapshot.Ports {
		if p.Protocol == row[0] && p.Address == row[1] && p.Remote == row[2] && strconv.Itoa(int(p.PID)) == row[4] {
			return p, true
		}
	}
	return scanner.PortEntry{}, false
}

func (m Model) historyView() string {
	if len(m.events) == 0 {
		return panelStyle.Width(m.width).Height(max(4, m.height-4)).Render("No port changes yet.")
	}
	lines := []string{}
	for _, event := range m.events {
		lines = append(lines, fmt.Sprintf("%s  %-7s  %-5s %-20s %s", event.At.Format("15:04:05"), event.Type, event.Port.Protocol, event.Port.Address, event.Port.Command))
	}
	return panelStyle.Width(m.width).Height(max(4, m.height-4)).Render(strings.Join(lines, "\n"))
}

func (m Model) confirmView() string {
	verb := "interrupt"
	if m.confirm.force {
		verb = "kill"
	}
	return panelStyle.Width(m.width).Height(max(4, m.height-4)).Render(fmt.Sprintf("Confirm %s pid %d (%s)?\n\nPress y to confirm, anything else to cancel.", verb, m.confirm.entry.PID, m.confirm.entry.Command))
}

func (m Model) footer() string {
	return "↑/↓ move  / filter  h history  r refresh  k interrupt  K kill  c copy  p port  o open  q quit"
}

func matches(p scanner.PortEntry, filter string) bool {
	if filter == "" {
		return true
	}
	haystack := strings.ToLower(strings.Join([]string{
		p.Protocol, p.Address, p.State, strconv.Itoa(int(p.PID)), p.Process, p.Command, p.Container, p.Image,
	}, " "))
	return strings.Contains(haystack, strings.ToLower(filter))
}

func rowSummary(p scanner.PortEntry) string {
	parts := []string{p.Protocol, p.Address, p.State, fmt.Sprintf("pid=%d", p.PID), p.Process, p.Command}
	if p.Container != "" {
		parts = append(parts, "container="+p.Container)
	}
	return strings.Join(parts, " ")
}

func report(err error, ok string) string {
	if err != nil {
		return err.Error()
	}
	return ok
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
