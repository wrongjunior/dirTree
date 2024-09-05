package tui

import (
	"fmt"
	"os"
	_ "path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"dirTree/internal/config"
	"dirTree/internal/scanner"
)

type model struct {
	config       config.Config
	filepicker   filepicker.Model
	progress     progress.Model
	viewport     viewport.Model
	state        state
	selectedDir  string
	scanResult   string
	scanError    error
	scannedFiles int
	totalFiles   int
}

type state int

const (
	stateFilePicker state = iota
	stateScanning
	stateResult
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func initialModel() model {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.Getwd()
	fp.AllowedTypes = []string{"dir"}

	return model{
		config:     config.Config{},
		filepicker: fp,
		progress:   progress.New(progress.WithDefaultGradient()),
		state:      stateFilePicker,
	}
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height-6)
		m.viewport.SetContent(m.scanResult)
		m.progress.Width = msg.Width - 4
		if m.state == stateResult {
			m.viewport.Height = msg.Height - 1
		}
	}

	switch m.state {
	case stateFilePicker:
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedDir = path
			m.state = stateScanning
			return m, m.startScanning
		}

	case stateScanning:
		switch msg := msg.(type) {
		case progressMsg:
			m.scannedFiles = msg.scanned
			m.totalFiles = msg.total
			if m.totalFiles > 0 {
				cmd = m.progress.SetPercent(float64(m.scannedFiles) / float64(m.totalFiles))
			}
		case scanCompleteMsg:
			m.scanResult = msg.result
			m.state = stateResult
			m.viewport.SetContent(m.scanResult)
		case scanErrorMsg:
			m.scanError = msg.err
			m.state = stateResult
		}

	case stateResult:
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateFilePicker:
		return titleStyle.Render("Select a directory to scan:") + "\n\n" + m.filepicker.View()
	case stateScanning:
		return fmt.Sprintf(
			"%s\n\n%s\n%d/%d files scanned",
			titleStyle.Render("Scanning directory:"),
			m.progress.View(),
			m.scannedFiles,
			m.totalFiles,
		)
	case stateResult:
		if m.scanError != nil {
			return errorStyle.Render(fmt.Sprintf("Error: %v", m.scanError))
		}
		return m.viewport.View()
	default:
		return "Unknown state"
	}
}

func (m model) startScanning() tea.Msg {
	m.config.OutputFile = ""
	m.config.RelativeFlag = true

	result, err := scanner.ScanDirectory(m.selectedDir, m.config, func(scanned, total int) {
		m.scannedFiles = scanned
		m.totalFiles = total
	})

	if err != nil {
		return scanErrorMsg{err}
	}

	return scanCompleteMsg{result}
}

type progressMsg struct {
	scanned, total int
}

type scanCompleteMsg struct {
	result string
}

type scanErrorMsg struct {
	err error
}

func Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
