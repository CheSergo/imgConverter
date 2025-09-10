package main

import (
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	dir      string
	input    string
	output   string
	fromType string
	toType   string
	depth    int
}

func (c *Config) CheckConfigType(t string) bool {
	suffix := strings.ToLower(t)
	suffix = strings.ReplaceAll(suffix, ".", "")
	return suffix == "jpg" || suffix == "jpeg" || suffix == "png" || suffix == "webp"
}

func EmptyConfig() *Config {
	return &Config{}
}

func NewConfig(input, output, dir, from, to string, depth int) (*Config, error) {
	if depth < 0 {
		return nil, fmt.Errorf("Depth cannot be less then 0")
	}

	if to == "jpg" {
		to = "jpeg"
	}

	return &Config{
		dir:      dir,
		input:    input,
		output:   output,
		fromType: from,
		toType:   to,
		depth:    depth,
	}, nil
}

func (c *Config) addDepth(depth int) error {
	if depth < 0 {
		return fmt.Errorf("Depth cannot be less then 0")
	}
	c.depth = depth
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
			c.fromType = typeToAdd
		case "to":
			c.toType = typeToAdd
		default:
			return errors.New("invalid target; use 'from' or 'to'")
		}
	} else {
		return fmt.Errorf("Type [%s] is not compare to jpeg|png|webp", target)
	}

	return nil
}
