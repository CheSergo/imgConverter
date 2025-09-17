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

// func hasField(walker DirectoryWalker, fieldName string) bool {
// 	val := reflect.ValueOf(cfg)
// 	field := val.FieldByName(fieldName)

// 	return field.IsValid()
// }

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

func (app *application) Walk() error {
	return filepath.WalkDir(app.walker.root, app.walkFunc)
}

func (app *application) walkFunc(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.Type()&fs.ModeSymlink != 0 {
		return nil
	}

	depth := strings.Count(path, string(os.PathSeparator)) - app.walker.baseDepth
	if depth > app.walker.maxDepth {
		return nil
	}
	app.logger.Info("Current path", "path", path)
	if depth == 0 {
		return nil
	}
	if !entry.IsDir() {
		if err := app.handleImageConversion(path, "", true); err != nil {
			return err
		}
	}
	return nil
}

func (app *application) processImage(path, newPath string, converter func(string, string) error) error {
	if len(newPath) == 0 {
		newPath = changeFileExtension(path, app.walker.toType)
	}

	err := converter(path, newPath)
	if err != nil {
		return fmt.Errorf("failed to convert %s to %s: %w", path, app.walker.toType, err)
	}
	return nil
}

func (app *application) handleImageConversion(path, newPath string, isFlaged bool) error {
	_, fileType, err := checkType(path)
	fmt.Println("Given fileType = ", fileType)
	if err != nil {
		return err
	}
	if !isFlaged || app.walker.fromType == fileType {
		switch fileType {
		case webpType:
			return app.processImage(path, newPath, convertWebP)
		case jpegType, jpgType:
			return app.processImage(path, newPath, convertJpeg)
		case pngType:
			return app.processImage(path, newPath, convertPng)
		}
	}
	return nil
}

func (app *application) processFile(path string) error {
	_, filetype, err := checkType(path)
	if err != nil {
		return err
	}

	if app.walker.fromType == filetype {
		switch filetype {
		case webpType:
			newPath := changeFileExtension(path, app.walker.toType)
			err := convertWebP(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, app.walker.toType, err)
			}
		case jpegType, jpgType:
			newPath := changeFileExtension(path, app.walker.toType)
			err := convertJpeg(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, app.walker.toType, err)
			}
		case pngType:
			newPath := changeFileExtension(path, app.walker.toType)
			err := convertPng(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, app.walker.toType, err)
			}
		default:
			return fmt.Errorf("unknown type")
		}
	}

	return nil
}
