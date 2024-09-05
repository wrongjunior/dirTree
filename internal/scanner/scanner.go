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

		if cfg.TUIMode && !shouldInclude(path, cfg.SelectedDirs) {
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

func shouldInclude(path string, selectedDirs []string) bool {
	if len(selectedDirs) == 0 {
		return true
	}
	for _, dir := range selectedDirs {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
}

// ScanDirectory scans the specified directory and returns the result as a string.
// It also takes a callback function to report scanning progress.
func ScanDirectory(dir string, cfg config.Config, progressCallback func(scanned, total int)) (string, error) {
	// Get the total number of files for progress tracking.
	totalFiles, err := countTotalFiles(dir, cfg)
	if err != nil {
		return "", fmt.Errorf("ошибка при подсчете файлов: %w", err)
	}

	// Initialize a slice to hold file information.
	var fileInfos []fileinfo.FileInfo

	// Walk through the directory structure and collect file information.
	scannedFiles := 0
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка доступа к пути %q: %w", path, err)
		}

		if shouldIgnore(path, info, cfg) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if cfg.TUIMode && !shouldInclude(path, cfg.SelectedDirs) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		fileInfos = append(fileInfos, fileinfo.FileInfo{Path: path, Info: info})
		scannedFiles++
		progressCallback(scannedFiles, totalFiles) // Call progress callback to update progress.

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("ошибка при обходе директории: %w", err)
	}

	// Create a writer to output the result.
	_, closer, err := output.GetWriter(cfg)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании writer: %w", err)
	}
	if closer != nil {
		defer closer.Close()
	}

	// Output the result based on the configuration (tree view or list).
	var resultBuilder strings.Builder
	if cfg.AbsoluteFlag || cfg.RelativeFlag {
		if err := output.OutputPathList(&resultBuilder, fileInfos, cfg, dir); err != nil {
			return "", err
		}
	} else {
		if err := output.OutputTreeView(&resultBuilder, fileInfos, cfg, dir); err != nil {
			return "", err
		}
	}

	// Return the result as a string.
	return resultBuilder.String(), nil
}

// countTotalFiles counts the total number of files in the directory for progress tracking.
func countTotalFiles(dir string, cfg config.Config) (int, error) {
	totalFiles := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка доступа к пути %q: %w", path, err)
		}

		if shouldIgnore(path, info, cfg) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if cfg.TUIMode && !shouldInclude(path, cfg.SelectedDirs) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		totalFiles++
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка при обходе директории для подсчета файлов: %w", err)
	}

	return totalFiles, nil
}
