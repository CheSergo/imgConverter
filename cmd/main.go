package main

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	dir      string
	input    string
	output   string
	fromType string
	toType   string
	depth    int
}

func main() {
	var cfg config

	flag.StringVar(&cfg.dir, "dir", "", "path to a dir")
	// flag.String("dir", "", "Путь к директории с изображениями")
	flag.StringVar(&cfg.fromType, "from", "", "Тип изображений для конвертации (например, jpg, png)")
	flag.StringVar(&cfg.toType, "to", "", "Тип изображений, в который нужно конвертировать (например, png, jpg)")
	flag.IntVar(&cfg.depth, "depth", 0, "Глубина обхода root")

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("No args")
		os.Exit(1)
	} else if len(args) == 2 {
		cfg.input = args[0]
		cfg.output = args[1]

		err := convertWebPToJPEGasChai(cfg.input, cfg.output)
		if err != nil {
			fmt.Println("Main error")
		}
	} else if len(args) >= 3 {
		flag.Parse()

		if cfg.dir != "" && cfg.fromType != "" && cfg.toType != "" {
			fmt.Printf("Checking dir: %v\n", cfg.dir)
			fmt.Println("Parsing files in the folder")
			var err error
			if cfg.depth != 0 {
				err, _ = walkDirectory(cfg.dir, cfg.depth)
			} else {
				err, _ = walkDirectory(cfg.dir, 1)
			}
			if err != nil {
				fmt.Printf("Error while walking: %v", err)
			}
		}
	}
	fmt.Println("Exit")
}
