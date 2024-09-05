package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	isSelected   bool
}

type state int

const (
	stateFilePicker state = iota
	stateScanning
	stateResult
)

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
)

func initialModel(cfg config.Config) model {
	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.Getwd()
	fp.AllowedTypes = []string{"dir"}

	return model{
		config:     cfg,
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

		case "enter":
			// Выбор директории с индикацией [x] при выборе
			if m.state == stateFilePicker && !m.isSelected {
				if m.filepicker.Path != "" {
					m.selectedDir = m.filepicker.Path
					m.isSelected = true
				}
			} else if m.isSelected && m.selectedDir != "" {
				// Если директория выбрана, начать сканирование
				m.state = stateScanning
				return m, tea.Batch(m.startScanning()) // Вызываем startScanning как tea.Cmd
			}

		case "space":
			// Отменить выбор при нажатии "space"
			if m.state == stateFilePicker && m.isSelected {
				m.isSelected = false
				m.selectedDir = ""
			}
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
		return titleStyle.Render("Выберите директорию для сканирования:") + "\n\n" +
			m.renderFilepickerView()
	case stateScanning:
		return fmt.Sprintf(
			"%s\n\n%s\n%d/%d файлов просканировано",
			titleStyle.Render("Сканирование директории:"),
			m.progress.View(),
			m.scannedFiles,
			m.totalFiles,
		)
	case stateResult:
		if m.scanError != nil {
			return errorStyle.Render(fmt.Sprintf("Ошибка: %v", m.scanError))
		}
		return m.viewport.View()
	default:
		return "Unknown state"
	}
}

func (m model) renderFilepickerView() string {
	// Обновляем отображение папок с индикацией выбора
	fpView := m.filepicker.View()
	fpView = m.addSelectionIndicators(fpView)
	return fpView
}

// Добавляем индикаторы [ ] и [x] в строки файлового пикера
func (m model) addSelectionIndicators(view string) string {
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		// Добавляем только для директорий
		if strings.HasPrefix(line, "d") {
			indicator := "[ ]"
			if m.isSelected && filepath.Base(m.selectedDir) == strings.TrimSpace(line) {
				indicator = "[x]"
			}
			lines[i] = fmt.Sprintf("%s %s", line, indicator)
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) startScanning() tea.Cmd {
	return func() tea.Msg {
		m.config.OutputFile = "" // Выводим результат на экран
		m.config.RelativeFlag = true

		// Сканируем директорию и выводим дерево директорий
		result, err := scanner.ScanDirectory(m.selectedDir, m.config, func(scanned, total int) {
			m.scannedFiles = scanned
			m.totalFiles = total
		})

		if err != nil {
			return scanErrorMsg{err}
		}

		return scanCompleteMsg{result}
	}
}

// Сообщения для обновления состояния
type progressMsg struct {
	scanned, total int
}

type scanCompleteMsg struct {
	result string
}

type scanErrorMsg struct {
	err error
}

// RunTUI запускает TUI и возвращает выбранные директории
func RunTUI(cfg config.Config) ([]string, error) {
	p := tea.NewProgram(initialModel(cfg), tea.WithAltScreen())
	finalModel, err := p.Run()

	if err != nil {
		return nil, fmt.Errorf("ошибка запуска TUI: %w", err)
	}

	if m, ok := finalModel.(model); ok && m.state == stateResult && m.scanError == nil {
		return []string{m.selectedDir}, nil
	}

	return nil, fmt.Errorf("TUI был закрыт без выбора директории")
}
