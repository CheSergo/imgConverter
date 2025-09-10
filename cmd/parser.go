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
	// fmt.Println("Current_depth is ", depth)
	// fmt.Println("Given_depth is ", dw.maxDepth)
	if depth > dw.maxDepth {
		// fmt.Printf("Current depth - %d is bigger then %d", depth, dw.maxDepth)
		// fmt.Println("Exit with if/else")
		return nil
	}
	fmt.Printf("Current path is - [%s]\n", path)
	if depth == 0 {
		fmt.Println("Skipping root dir")
		return nil
	}
	// fmt.Println("We are here")
	if !entry.IsDir() {
		// fmt.Println("not a dir")
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
			fmt.Println("Формат: WEBp")
			newPath := changeFileExtension(path, dw.toType)
			err := convertWebP(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, dw.toType, err)
			}
		case jpegType, jpgType:
			fmt.Println("Формат: JPEG")
			newPath := changeFileExtension(path, dw.toType)
			err := convertJpeg(path, newPath)
			if err != nil {
				return fmt.Errorf("failed to convert %s to %s: %w", path, dw.toType, err)
			}
		case pngType:
			fmt.Println("Формат: PNG")
		default:
			fmt.Println("Неизвестный формат")
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

	// if filepath.IsDir() {
	// 	return true, nil
	// } else {
	// 	return false, nil
	// }
}

// func walkDirectory(root string, maxDepth int) (error, []string) {
// 	absRoot, err := filepath.Abs(root)
// 	if err != nil {
// 		return err, nil
// 	}

// 	baseDepth := strings.Count(absRoot, string(os.PathSeparator))
// 	var list []string
// 	firstIteration := true

// 	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if entry.Type()&fs.ModeSymlink != 0 {
// 			return nil
// 		}

// 		if firstIteration {
// 			firstIteration = false
// 			b, err := checkIsDir(path)
// 			if err != nil {
// 				return err
// 			}
// 			if b {
// 				return nil
// 			} else {
// 				return fmt.Errorf("Path is a file")
// 			}
// 		}

// 		absRoot, err := filepath.Abs(path)
// 		if err != nil {
// 			return nil
// 		}

// 		depth := strings.Count(absRoot, string(os.PathSeparator)) - baseDepth
// 		if depth <= maxDepth {
// 			if entry.IsDir() {
// 				return fs.SkipDir
// 			} else {
// 				_, tp, err := checkType(path)
// 				if err != nil {
// 					return err
// 				}

// 				types := formatsList()
// 				switch tp {
// 				case types.webp:
// 					newPath := changeFileExtension(path, "jpg")
// 					err := convertWebPToJPEGasChai(path, newPath)
// 					if err != nil {
// 						return fmt.Errorf("failed to convert %s to JPG: %w", path, err)
// 					}
// 				case types.jpeg:
// 					fmt.Println("Формат: JPEG")
// 				case types.png:
// 					fmt.Println("Формат: PNG")
// 				default:
// 					fmt.Println("Неизвестный формат")
// 				}
// 			}
// 			return nil
// 		}
// 		return nil
// 	})

// 	return err, list
// }
