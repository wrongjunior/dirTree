package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dirTree/internal/config"
	"dirTree/internal/fileinfo"
	"dirTree/internal/output"
)

func Run(cfg config.Config) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ошибка при получении текущей директории: %w", err)
	}

	fileInfos, err := walkDirectory(currentDir, cfg)
	if err != nil {
		return fmt.Errorf("ошибка при обходе директории: %w", err)
	}

	writer, closer, err := output.GetWriter(cfg)
	if err != nil {
		return fmt.Errorf("ошибка при создании writer: %w", err)
	}
	if closer != nil {
		defer func() {
			if err := closer.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при закрытии writer: %v\n", err)
			}
		}()
	}

	if cfg.AbsoluteFlag || cfg.RelativeFlag {
		if err := output.OutputPathList(writer, fileInfos, cfg, currentDir); err != nil {
			return err
		}
	} else {
		if err := output.OutputTreeView(writer, fileInfos, cfg, currentDir); err != nil {
			return err
		}
	}

	return nil
}

func walkDirectory(root string, cfg config.Config) ([]fileinfo.FileInfo, error) {
	var fileInfos []fileinfo.FileInfo
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка доступа к пути %q: %w", path, err)
		}
		if shouldIgnore(path, info, cfg) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		fileInfos = append(fileInfos, fileinfo.FileInfo{Path: path, Info: info})
		return nil
	})
	return fileInfos, err
}

func shouldIgnore(path string, info os.FileInfo, cfg config.Config) bool {
	if info.IsDir() {
		for _, dir := range cfg.IgnoreDirs {
			if strings.HasSuffix(path, dir) {
				return true
			}
		}
	} else {
		if cfg.IgnoreExts[filepath.Ext(path)] {
			return true
		}
	}
	return false
}
