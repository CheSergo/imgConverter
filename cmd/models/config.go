package models

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Dir      string
	Input    string
	Output   string
	FromType string
	ToType   string
	Depth    int
}

func (c *Config) CheckConfigType(t string) bool {
	suffix := strings.ToLower(t)
	suffix = strings.ReplaceAll(suffix, ".", "")
	return suffix == "jpg" || suffix == "jpeg" || suffix == "png" || suffix == "webp"
}

func EmptyConfig() *Config {
	return &Config{}
}

func NewConfig(args []string) *Config {
	cfg := EmptyConfig()

	if len(args) == 2 {
		cfg.Input = args[0]
		cfg.Output = args[1]
	} else {
		flag.StringVar(&cfg.Dir, "dir", "", "Path to the directory with images")
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
	}

	return cfg
}

func (c *Config) ValidateConfig() error {
	if c.Dir == "" || c.FromType == "" || c.ToType == "" {
		return fmt.Errorf("Error: Directory, fromType, and toType must be specified.")
	}

	if c.FromType == c.ToType {
		return fmt.Errorf("Error with types. Convert types must be different.")
	}

	return nil
}

func (c *Config) addDepth(depth int) error {
	if depth < 0 {
		return fmt.Errorf("Depth cannot be less then 0")
	}
	c.Depth = depth
	return nil
}

func (c *Config) setType(t, target string) error {
	typeToAdd := strings.ToLower(t)
	typeToAdd = strings.ReplaceAll(typeToAdd, ".", "")

	if c.CheckConfigType(typeToAdd) {
		if typeToAdd == "jpg" {
			typeToAdd = "jpeg"
		}

		switch target {
		case "from":
			c.FromType = typeToAdd
		case "to":
			c.ToType = typeToAdd
		default:
			return errors.New("invalid target; use 'from' or 'to'")
		}
	} else {
		return fmt.Errorf("Type [%s] is not compare to jpeg|png|webp", target)
	}

	return nil
}
