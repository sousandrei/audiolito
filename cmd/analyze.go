package main

import (
	"fmt"
	"normalizer/internal/normalizer"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func handleAnalyze(c *cli.Context) error {
	ctx := c.Context

	noTui := c.Bool("no-tui")
	isDebug := c.Bool("debug")

	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		stats, err := normalizer.Analyze(
			ctx,
			normalizer.WithFilePath(inputFilePath),
			normalizer.DisableTui(noTui),
			normalizer.WithDebug(isDebug),
		)
		if err != nil {
			return err
		}

		fmt.Printf("       file: %s\n", inputFilePath)
		fmt.Printf("mean_volume: %s\n", strconv.FormatFloat(stats.Mean_volume, 'f', -1, 64))
		fmt.Printf(" max_volume: %s\n\n", strconv.FormatFloat(stats.Max_volume, 'f', -1, 64))
	}
	return nil
}
