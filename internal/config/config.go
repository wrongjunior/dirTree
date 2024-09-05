package config

import (
	"bufio"
	"flag"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	OutputFile       string
	RelativeFlag     bool
	AbsoluteFlag     bool
	IgnoreDirs       []string
	IgnoreExts       map[string]bool
	IgnoreConfigFile string
	ClipboardFlag    bool
	TUIMode          bool
	SelectedDirs     []string
}

func ParseFlags() Config {
	relativeFlag := flag.Bool("relative", false, "Выводить только относительные пути")
	absoluteFlag := flag.Bool("absolute", false, "Выводить только абсолютные пути")
	outputFileFlag := flag.String("output", "", "Имя файла для записи вывода")
	ignoreDirsFlag := flag.String("ignore-dirs", "", "Список игнорируемых директорий (через запятую)")
	ignoreExtsFlag := flag.String("ignore-exts", "", "Список игнорируемых расширений файлов (через запятую)")
	ignoreConfigFileFlag := flag.String("ignore-config", "", "Файл конфигурации с игнорируемыми директориями и расширениями")
	clipboardFlag := flag.Bool("clipboard", false, "Выводить в буфер обмена")
	tuiFlag := flag.Bool("tui", false, "Запустить в интерактивном режиме выбора папок")

	flag.Parse()

	cfg := Config{
		OutputFile:       *outputFileFlag,
		RelativeFlag:     *relativeFlag,
		AbsoluteFlag:     *absoluteFlag,
		IgnoreDirs:       filterEmpty(strings.Split(*ignoreDirsFlag, ",")),
		IgnoreExts:       convertToExtMap(filterEmpty(strings.Split(*ignoreExtsFlag, ","))),
		IgnoreConfigFile: *ignoreConfigFileFlag,
		ClipboardFlag:    *clipboardFlag,
		TUIMode:          *tuiFlag,
	}

	if cfg.IgnoreConfigFile != "" {
		loadIgnoreConfig(&cfg)
	}

	return cfg
}

func filterEmpty(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func convertToExtMap(exts []string) map[string]bool {
	extMap := make(map[string]bool)
	for _, ext := range exts {
		extMap["."+ext] = true
	}
	return extMap
}

func loadIgnoreConfig(cfg *Config) {
	file, err := os.Open(cfg.IgnoreConfigFile)
	if err != nil {
		slog.Error("Ошибка при открытии файла конфигурации", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Ошибка при закрытии файла конфигурации", slog.String("error", err.Error()))
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "dir:") {
			cfg.IgnoreDirs = append(cfg.IgnoreDirs, strings.TrimPrefix(line, "dir:"))
		} else if strings.HasPrefix(line, "ext:") {
			cfg.IgnoreExts["."+strings.TrimPrefix(line, "ext:")] = true
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Ошибка при чтении файла конфигурации", slog.String("error", err.Error()))
	}
}
