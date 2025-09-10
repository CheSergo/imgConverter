package main

import "strings"

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
