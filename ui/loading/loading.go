package loading

import (
	"context"
	"os"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Loading struct {
	LoadingModel spinner.Model
	Quitting     chan bool
	ParentCtx    context.Context
	Ctx          context.Context
	Can          context.CancelFunc
	Wg           sync.WaitGroup
	Messgage     string
}

func NewLoading(ParentCtx context.Context) *Loading {
	ctx, can := context.WithCancel(context.Background())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Loading{
		LoadingModel: s,
		Quitting:     make(chan bool),
		ParentCtx:    ParentCtx,
		Ctx:          ctx,
		Can:          can,
		Wg:           sync.WaitGroup{},
		Messgage:     "Calculating template count...",
	}
}

func (m *Loading) Run() {
	m.Wg.Add(1)
	go func() {
		defer m.Can()
		defer m.Wg.Done()
		tea.NewProgram(m, tea.WithContext(m.ParentCtx)).Run()
	}()
}

func (m *Loading) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, m.LoadingModel.Tick)
}

func (m *Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			close(m.Quitting)
			return m, tea.Quit
		default:
			return m, nil
		}

	default:
		select {
		case <-m.Quitting:
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.LoadingModel, cmd = m.LoadingModel.Update(msg)
			return m, cmd
		}
	}
}

func (m *Loading) View() string {

	str := "\n " + m.LoadingModel.View() + " " + m.Messgage + "\n\n"
	select {
	case <-m.Quitting:
		return ""
	default:
		return str + "\n"
	}
}

func (m *Loading) End() {
	select {
	case <-m.Quitting:
		os.Exit(1)
		return
	default:
		close(m.Quitting)
	}
}
