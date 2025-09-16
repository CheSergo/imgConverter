package main

import (
	"imgConverter/cmd/models"
	"log/slog"
	"os"
)

type application struct {
	logger *slog.Logger
	// config *models.Config
	walker *DirectoryWalker
}

func main() {
	// init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	args := os.Args[1:]
	if len(args) < 2 {
		app.logger.Error("Error: At least two arguments are required.")
		return
	}

	cfg := models.NewConfig(args)

	if len(args) == 2 {
		if err := convertWebP(cfg.Input, cfg.Output); err != nil {
			app.logger.Error("Error while converting files", "error", err.Error())
		}
		return
	}

	if err := cfg.ValidateConfig(); err != nil {
		app.logger.Error("Valodation error.", "Err", err)
	}

	app.logger.Info("Checking directory", "dir", cfg.Dir)
	app.logger.Info("Parsing files in the folder")

	path, err := os.Stat(cfg.Dir)
	if err != nil {
		app.logger.Error("Error while cheking the path", "path", path, "err", err)
	}
	if !path.IsDir() {
		app.logger.Error("Error. Given path is not a dir")
		return
	}
	if cfg.Depth == 0 {
		cfg.Depth = 1
	}

	app.walker, err = NewDirectoryWalker(cfg.Dir, cfg.FromType, cfg.ToType, cfg.Depth)
	if err != nil {
		app.logger.Error("Error creating new directory walker", "error", err)
	}
	if err := app.Walk(); err != nil {
		app.logger.Error("Error while walking the directory", "error", err)
	}

	app.logger.Info("Exit")
}
