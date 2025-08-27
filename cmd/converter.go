package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	chai "github.com/chai2010/webp"
	"golang.org/x/image/webp"
)

const (
	Png  = "png"
	Jpeg = "jpeg"
	Webp = "webp"
)

type Formats struct {
	png  string
	jpeg string
	webp string
}

func formatsList() *Formats {
	return &Formats{
		png:  Png,
		jpeg: Jpeg,
		webp: Webp,
	}
}

func convertWebPToJPEG(inputPath, outputPath string) error {
	webpFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Error while opening file: %v", err)
	}
	defer webpFile.Close()

	img, err := webp.Decode(webpFile)
	if err != nil {
		return fmt.Errorf("Error while decoding file WebP: %v", err)
	}

	jpegFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Error while creating file JPEG: %v", err)
	}
	defer jpegFile.Close()

	err = jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return fmt.Errorf("Error while converting the file to WebP: %v", err)
	}
	return nil
}

func convertWebPToJPEGasChai(inputPath, outputPath string) error {
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

	// GetMetadata
	// if metadata, err := chai.GetMetadata(file, "ICCP"); err != nil {
	// 	fmt.Printf("Metadata: err = %v\n", err)
	// } else {
	// 	fmt.Printf("Metadata: %s\n", string(metadata))
	// }

	m, err := chai.Decode(bytes.NewReader(file))
	if err != nil {
		return err
	}

	jpegFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	err = jpeg.Encode(jpegFile, m, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	fmt.Printf("Save the file - [%s]\n", outputPath)

	// Decode
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
