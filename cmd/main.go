package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	cfg := EmptyConfig()

	flag.StringVar(&cfg.dir, "dir", "", "Path to the directory with images")
	fromType := flag.String("from", "", "Type of images to convert (e.g. jpg, png)")
	toType := flag.String("to", "", "The type of images to convert to (e.g. png, jpg)")
	dep := flag.Int("depth", 1, "Root folder traversal depth")
	flag.Parse()

	err := cfg.setType(*fromType, "from")
	if err != nil {
		fmt.Println(err)
	}
	err = cfg.setType(*toType, "to")
	if err != nil {
		fmt.Println(err)
	}
	cfg.addDepth(*dep)

	// Checking amount of args
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("No args")
		os.Exit(1)
	} else if len(args) == 2 {
		cfg.input = args[0]
		cfg.output = args[1]

		err := convertWebP(cfg.input, cfg.output)
		if err != nil {
			fmt.Println("Error while converting files")
		}
	} else if len(args) >= 3 {

		if cfg.dir != "" && cfg.fromType != "" && cfg.toType != "" {

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

			// Main logic
			walker, err := NewDirectoryWalker(cfg.dir, cfg.fromType, cfg.toType, cfg.depth)
			if err != nil {
				fmt.Printf("Error at new directory walker func")
			}
			fmt.Printf("New obj - %+v\n", walker)
			if err := walker.Walk(); err != nil {
				fmt.Println("Ошибка при обходе директории:", err)
			}
		}
	}
	fmt.Println("Exit")
}
