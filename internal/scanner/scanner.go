package scanner

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dirTree/internal/config"
	"dirTree/internal/fileinfo"
	"dirTree/internal/output"
)

func ScanDirectory(root string, cfg config.Config, progressCallback func(int, int)) (string, error) {
	fileInfos, err := walkDirectory(root, cfg, progressCallback)
	if err != nil {
		return "", fmt.Errorf("ошибка при обходе директории: %w", err)
	}

	var buf bytes.Buffer
	writer := output.GetBufferWriter(&buf)

	if cfg.AbsoluteFlag || cfg.RelativeFlag {
		if err := output.OutputPathList(writer, fileInfos, cfg, root); err != nil {
			return "", err
		}
	} else {
		if err := output.OutputTreeView(writer, fileInfos, cfg, root); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func walkDirectory(root string, cfg config.Config, progressCallback func(int, int)) ([]fileinfo.FileInfo, error) {
	var fileInfos []fileinfo.FileInfo
	var totalFiles int

	// First pass: count total files
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("ошибка доступа к пути %q: %w", path, err)
		}
		if !shouldIgnore(path, info, cfg) {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Second pass: collect file info and update progress
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		progressCallback(len(fileInfos), totalFiles)
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
