// main.go
package main

import (
	_ "fmt"
	"log/slog"
	"os"

	"dirTree/internal/config"
	"dirTree/internal/scanner"
	"dirTree/internal/tui"
)

func main() {
	cfg := config.ParseFlags()

	if cfg.TUIMode {
		// Передаём cfg в RunTUI
		selectedDirs, err := tui.RunTUI(cfg)
		if err != nil {
			slog.Error("Ошибка в интерактивном режиме", slog.String("error", err.Error()))
			os.Exit(1)
		}
		cfg.SelectedDirs = selectedDirs
	}

	if cfg.RelativeFlag && cfg.AbsoluteFlag {
		slog.Error("Взаимоисключающие флаги: --relative и --absolute не могут быть установлены одновременно.")
		os.Exit(1)
	}

	if err := scanner.Run(cfg); err != nil {
		slog.Error("Ошибка", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
