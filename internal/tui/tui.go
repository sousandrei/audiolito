package tui

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type ProgressMsg float64
type SetDurationMsg float64
type FinishMsg bool

type Model struct {
	progress progress.Model
	duration float64
}

func New(ctx context.Context) (*tea.Program, string, error) {
	p := tea.NewProgram(Model{
		progress: progress.New(progress.WithDefaultGradient()),
	})

	ps, err := newServer(ctx, p)
	if err != nil {
		return nil, "", err
	}

	return p, ps.addr, nil
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return m.progress.SetPercent(0.0)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case SetDurationMsg:
		m.duration = float64(msg)
		return m, nil

	case ProgressMsg:
		// m.duration is total, progressMsg is current
		percent := float64(msg) / m.duration

		if percent >= 1.0 {
			return m, m.progress.SetPercent(1.0)
		}

		return m, m.progress.SetPercent(percent)

	default:
		return m, nil
	}
}

func (m Model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Press q or ctrl+c to quit")
}
