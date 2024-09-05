package main

import (
	"log/slog"
	"os"

	"dirTree/internal/config"
	"dirTree/internal/scanner"
)

func main() {
	cfg := config.ParseFlags()

	if cfg.RelativeFlag && cfg.AbsoluteFlag {
		slog.Error("Взаимоисключающие флаги: --relative и --absolute не могут быть установлены одновременно.")
		os.Exit(1)
	}

	if err := scanner.Run(cfg); err != nil {
		slog.Error("Ошибка", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
