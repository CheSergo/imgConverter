package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	webpType = "webp"
	jpegType = "jpeg"
	jpgType  = "jpg"
	pngType  = "png"
)

type DirectoryWalker struct {
	root       string
	maxDepth   int
	baseDepth  int
	fromType   string
	toType     string
	foundFiles []string
}

func NewDirectoryWalker(root, from, to string, maxDepth int) (*DirectoryWalker, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	return &DirectoryWalker{
		root:      absRoot,
		maxDepth:  maxDepth,
		baseDepth: strings.Count(absRoot, string(os.PathSeparator)),
		fromType:  from,
		toType:    to,
	}, nil
}

func (dw *DirectoryWalker) Walk() error {
	return filepath.WalkDir(dw.root, dw.walkFunc)
}

func (dw *DirectoryWalker) walkFunc(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.Type()&fs.ModeSymlink != 0 {
		return nil
	}

	depth := strings.Count(path, string(os.PathSeparator)) - dw.baseDepth
	if depth > dw.maxDepth {
		return nil
	}
	fmt.Printf("Current path is - [%s]\n", path)
	if depth == 0 {
		fmt.Println("Skipping root dir")
		return nil
	}
	if !entry.IsDir() {
		if err := dw.processFile(path); err != nil {
			return err
		}
	}
	return nil
}

func (dw *DirectoryWalker) processFile(path string) error {
	_, filetype, err := checkType(path)
	if err != nil {
		return err
	}

	if dw.fromType == filetype {
		switch filetype {
		case webpType:
			// fmt.Println("Формат: WEBp")
			newPath := changeFileExtension(path, dw.toType)
			err := convertWebP(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, dw.toType, err)
			}
		case jpegType, jpgType:
			// fmt.Println("Формат: JPEG")
			newPath := changeFileExtension(path, dw.toType)
			err := convertJpeg(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, dw.toType, err)
			}
		case pngType:
			// fmt.Println("Формат: PNG")
			newPath := changeFileExtension(path, dw.fromType)
			err := convertPng(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, dw.toType, err)
			}
		default:
			fmt.Println("unknown type")
		}
	}

	return nil
}

func checkPath(path string) error {
	filepath, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Error while checking the path %s: %v", path, err)
	}

	if filepath.IsDir() {
		fmt.Printf("This is a folder: %s", path)
	} else {
		fmt.Printf("This is a file: %s", path)
	}

	return nil
}

func checkIsDir(path string) (bool, error) {
	filepath, err := os.Stat(path)
	if err != nil {
		err = fmt.Errorf("Error while checking the path %s: %v", path, err)
		return false, err
	}

	return filepath.IsDir(), nil

}
