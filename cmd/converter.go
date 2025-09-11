package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	chai "github.com/chai2010/webp"
)

const (
	Png  = "png"
	Jpeg = "jpeg"
	Webp = "webp"
)

func convertWebP(inputPath, outputPath string) error {
	var width, height int

	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// GetInfo
	if width, height, _, err = chai.GetInfo(file); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("width = %d, height = %d\n", width, height)

	m, err := chai.Decode(bytes.NewReader(file))
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	isPNG := strings.HasSuffix(outputPath, ".png")
	if isPNG {
		err = png.Encode(outputFile, m)
	} else {
		err = jpeg.Encode(outputFile, m, &jpeg.Options{Quality: 100})
	}

	if err != nil {
		return err
	}
	fmt.Printf("Save the file - [%s]\n", outputPath)

	// Decode
	return nil
}

func convertJpeg(inputPath, outputPath string) error {
	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	jpegFile, err := jpeg.Decode(bytes.NewReader(file))
	if err != nil {
		return err
	}

	// Creatig file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	isPNG := strings.HasSuffix(outputPath, ".png")
	if isPNG {
		err = png.Encode(outputFile, jpegFile)
	} else {
		err = chai.Encode(outputFile, jpegFile, &chai.Options{Lossless: true})
	}
	if err != nil {
		return err
	}

	return nil
}

func convertPng(inputPath, outputPath string) error {
	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	pngFile, err := png.Decode(bytes.NewBuffer(file))
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	isWebp := strings.HasSuffix(outputPath, ".webp")
	if isWebp {
		err = chai.Encode(outputFile, pngFile, &chai.Options{Lossless: true})
	} else {
		err = jpeg.Encode(outputFile, pngFile, &jpeg.Options{Quality: 100})
	}

	if err != nil {
		return err
	}

	return nil
}

func checkType(path string) (image.Image, string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("Error while opening the file: [%s]", path)
	}

	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("Error while deconding format: %v", err)
	}
	// fmt.Printf("File format - %s\n", format)

	return img, format, nil
}

func changeFileExtension(path, newExt string) string {
	// Получаем базовое имя файла без расширения
	base := filepath.Base(path)
	fmt.Printf("Base - [%s]\n", base)
	// Получаем имя файла без расширения
	name := base[:len(base)-len(filepath.Ext(base))]
	fmt.Printf("Name - [%s]\n", name)
	// Формируем новый путь с новым расширением
	newPath := filepath.Join(filepath.Dir(path), name+"."+newExt)
	fmt.Printf("Newpath - [%s]\n", newPath)
	return newPath
}
