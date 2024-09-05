package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"dirTree/internal/config"
	"dirTree/internal/fileinfo"
	"github.com/atotto/clipboard"
)

func GetWriter(cfg config.Config) (*bufio.Writer, io.Closer, error) {
	if cfg.ClipboardFlag {
		var buf strings.Builder
		return bufio.NewWriter(&buf), &bufferCloser{&buf}, nil
	}
	if cfg.OutputFile == "" {
		return bufio.NewWriter(os.Stdout), nil, nil
	}
	file, err := os.Create(cfg.OutputFile)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при создании файла: %w", err)
	}
	return bufio.NewWriter(file), file, nil
}

type bufferCloser struct {
	buf *strings.Builder
}

func (b *bufferCloser) Close() error {
	return clipboard.WriteAll(b.buf.String())
}

func OutputPathList(writer *bufio.Writer, fileInfos []fileinfo.FileInfo, cfg config.Config, rootDir string) error {
	for _, info := range fileInfos {
		displayPath := GetDisplayPath(info.Path, rootDir, cfg)
		if _, err := fmt.Fprintf(writer, "%s\n", displayPath); err != nil {
			return err
		}
	}
	return writer.Flush()
}

func OutputTreeView(writer *bufio.Writer, fileInfos []fileinfo.FileInfo, cfg config.Config, rootDir string) error {
	tree := make(map[string][]fileinfo.FileInfo)
	var roots []fileinfo.FileInfo

	for _, info := range fileInfos {
		dir := filepath.Dir(info.Path)
		relDir, err := filepath.Rel(rootDir, dir)
		if err != nil {
			return err
		}
		relDir = filepath.ToSlash(relDir)
		if relDir == "." {
			roots = append(roots, info)
		} else {
			tree[relDir] = append(tree[relDir], info)
		}
	}

	if _, err := fmt.Fprintf(writer, "%s/\n", filepath.Base(rootDir)); err != nil {
		return err
	}
	for i, root := range roots {
		prefix := "├── "
		lastPrefix := "│   "
		if i == len(roots)-1 {
			prefix = "└── "
			lastPrefix = "    "
		}
		if err := printTree(writer, root, prefix, lastPrefix, tree, cfg, rootDir); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func printTree(writer *bufio.Writer, item fileinfo.FileInfo, prefix, lastPrefix string, tree map[string][]fileinfo.FileInfo, cfg config.Config, rootDir string) error {
	displayPath := GetDisplayPath(item.Path, rootDir, cfg)

	if item.Info.IsDir() {
		if _, err := fmt.Fprintf(writer, "%s%s/\n", prefix, filepath.Base(displayPath)); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(writer, "%s%s (%s)\n", prefix, filepath.Base(displayPath), FormatFileSize(item.Info.Size())); err != nil {
			return err
		}
	}

	relPath, err := filepath.Rel(rootDir, item.Path)
	if err != nil {
		return err
	}
	relPath = filepath.ToSlash(relPath)

	children := tree[relPath]
	for i, child := range children {
		newPrefix := lastPrefix
		newLastPrefix := lastPrefix
		if i == len(children)-1 {
			newPrefix += "└── "
			newLastPrefix += "    "
		} else {
			newPrefix += "├── "
			newLastPrefix += "│   "
		}
		if err := printTree(writer, child, newPrefix, newLastPrefix, tree, cfg, rootDir); err != nil {
			return err
		}
	}
	return nil
}

func GetDisplayPath(path, rootDir string, cfg config.Config) string {
	if cfg.AbsoluteFlag {
		return path
	}
	if cfg.RelativeFlag {
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return path
		}
		return filepath.ToSlash(relPath)
	}
	return path
}

func FormatFileSize(size int64) string {
	switch {
	case size < fileinfo.BytesInKB:
		return fmt.Sprintf("%d B", size)
	case size < fileinfo.BytesInMB:
		return fmt.Sprintf("%.2f KB", float64(size)/fileinfo.BytesInKB)
	case size < fileinfo.BytesInGB:
		return fmt.Sprintf("%.2f MB", float64(size)/fileinfo.BytesInMB)
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/fileinfo.BytesInGB)
	}
}
