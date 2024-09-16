package progress

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

type (
	Progress struct {
		ProgressModel progress.Model
		PercentChan   chan float64
		IsComplete    float64
		ParentCtx     context.Context
		Ctx           context.Context
		Can           context.CancelFunc
		Wg            sync.WaitGroup
		TimeDuration  time.Duration
		Count         int
		Step          float64
	}

	tickMsg time.Time
)

func NewProgress(parentCtx context.Context) *Progress {
	ctx, can := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	return &Progress{
		ProgressModel: progress.New(progress.WithDefaultGradient()),
		PercentChan:   make(chan float64, 1),
		IsComplete:    1.0,
		ParentCtx:     parentCtx,
		Ctx:           ctx,
		Can:           can,
		Wg:            sync.WaitGroup{},
		TimeDuration:  50 * time.Millisecond,
		Count:         100,
	}
}

func (m *Progress) Run() {
	m.ConfigDuration()
	m.Wg.Add(1)
	go func() {
		defer m.Can()
		defer m.Wg.Done()
		tea.NewProgram(m, tea.WithContext(m.ParentCtx)).Run()
	}()
}

func (m *Progress) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tickCmd(m.TimeDuration))
}

func (m *Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.ProgressModel.Width = msg.Width - padding*2 - 4
		if m.ProgressModel.Width > maxWidth {
			m.ProgressModel.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		ProgressModel, cmd := m.ProgressModel.Update(msg)
		m.ProgressModel = ProgressModel.(progress.Model)
		return m, cmd

	case tickMsg:
		if m.ProgressModel.Percent() >= m.IsComplete {
			return m, tea.Tick(500*time.Millisecond, tickCallback(tea.Quit))
		}
		select {
		case percent := <-m.PercentChan:
			return m, tea.Batch(tickCmd(m.TimeDuration), m.ProgressModel.IncrPercent(percent))
		default:
			return m, tickCmd(m.TimeDuration)
		}

	default:
		return m, nil
	}
}

func (m *Progress) View() string {
	pad := strings.Repeat(" ", padding)

	return "\n" + pad + m.ProgressModel.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func (m *Progress) ConfigDuration() {
	m.Step = 1.0 / float64(m.Count)
	m.TimeDuration = time.Duration(int(m.Step*1000)) * time.Millisecond
}

func (p *Progress) UpdateProgress() {
	select {
	case <-p.Ctx.Done():
		if p.ProgressModel.Percent() < p.IsComplete {
			os.Exit(1)
		}
		return
	case p.PercentChan <- p.Step:
		return
	}
}

func tickCallback(t tea.Msg) func(time.Time) tea.Msg {
	return func(time time.Time) tea.Msg {
		tType := reflect.TypeOf(t)
		if tType.Kind() == reflect.Func {
			results := reflect.ValueOf(t).Call(nil)
			if len(results) > 0 {
				return results[0].Interface()
			}
			return nil
		}
		return t
	}
}

func tickCmd(td time.Duration) tea.Cmd {
	return tea.Tick(td, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

