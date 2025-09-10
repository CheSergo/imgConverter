package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var cfg Config

	flag.StringVar(&cfg.dir, "dir", "", "path to a dir")
	// flag.String("dir", "", "Путь к директории с изображениями")
	flag.StringVar(&cfg.fromType, "from", "", "Тип изображений для конвертации (например, jpg, png)")
	flag.StringVar(&cfg.toType, "to", "", "Тип изображений, в который нужно конвертировать (например, png, jpg)")
	flag.IntVar(&cfg.depth, "depth", 1, "Глубина обхода root")

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
			if !cfg.CheckConfigType(cfg.fromType) || !cfg.CheckConfigType(cfg.toType) {
				fmt.Println("Suffix is not compare to jpg|png|webp")
				return
			}

			if cfg.fromType == cfg.toType {
				fmt.Println("Error with types. Convert types must be different")
				return
			}

			fmt.Printf("Checking dir: %v\n", cfg.dir)
			fmt.Println("Parsing files in the folder")

			path, err := os.Stat(cfg.dir)
			if err != nil {
				fmt.Printf("Error while cheking the path %s: %v\n", path, err)
			}
			if !path.IsDir() {
				fmt.Println("Error. Given path is not a dir")
				return
			}

			if cfg.depth == 0 {
				cfg.depth = 1
			}
			walker, err := NewDirectoryWalker(cfg.dir, cfg.fromType, cfg.toType, cfg.depth)
			if err != nil {
				fmt.Printf("Error at new directory walker func")
			}
			fmt.Printf("New obj - %+v\n", walker)
			if err := walker.Walk(); err != nil {
				fmt.Println("Ошибка при обходе директории:", err)
			}
			// fmt.Printf("root - %v\n", walker.root)
			// fmt.Printf("Depth - %v\n", walker.maxDepth)
			// fmt.Printf("base - %v\n", walker.baseDepth)
			// fmt.Printf("list - %v\n", walker.foundFiles)

			// var err error
			// if cfg.depth != 0 {
			// 	err, _ = walkDirectory(cfg.dir, cfg.depth)
			// } else {
			// 	err, _ = walkDirectory(cfg.dir, 1)
			// }
			// if err != nil {
			// 	fmt.Printf("Error while walking: %v", err)
			// }
		}
	}
	fmt.Println("Exit")
}

// walker, err := NewDirectoryWalker("/path/to/directory", 3)
//     if err != nil {
//         log.Fatal(err)
//     }

//     if err := walker.Walk(); err != nil {
//         log.Println("Ошибка при обходе директории:", err)
//     }
