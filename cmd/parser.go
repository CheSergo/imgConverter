package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

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

	if filepath.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func walkDirectory(root string, maxDepth int) (error, []string) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err, nil
	}

	baseDepth := strings.Count(absRoot, string(os.PathSeparator))
	var list []string
	firstIteration := true

	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if firstIteration {
			firstIteration = false
			b, err := checkIsDir(path)
			if err != nil {
				return err
			}
			if b {
				return nil
			} else {
				return fmt.Errorf("Path is a file")
			}
		}

		absRoot, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		depth := strings.Count(absRoot, string(os.PathSeparator)) - baseDepth
		if depth <= maxDepth {
			if entry.IsDir() {
				return fs.SkipDir
			} else {
				_, tp, err := checkType(path)
				if err != nil {
					return err
				}

				types := formatsList()
				switch tp {
				case types.webp:
					newPath := changeFileExtension(path, "jpg")
					err := convertWebPToJPEGasChai(path, newPath)
					if err != nil {
						return err
					}
				case types.jpeg:
					fmt.Println("Формат: JPEG")
				case types.png:
					fmt.Println("Формат: PNG")
				default:
					fmt.Println("Неизвестный формат")
				}
			}
			return nil
		}
		return nil
	})

	return err, list
}
