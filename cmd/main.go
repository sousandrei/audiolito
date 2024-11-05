package main

import (
	"fmt"
	"log"
	"normalizer/internal/normalizer"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "audiolito",
		Usage: "tools for audio manipulation",
		Commands: []*cli.Command{
			{
				Name:      "normalize",
				ArgsUsage: "[FILES]",
				Action:    handleNormalize,
			},
			{
				Name:      "analyze",
				ArgsUsage: "[FILES]",
				Action:    handleAnalyze,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func handleNormalize(c *cli.Context) error {
	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		fmt.Printf("normalizing %s\n", inputFilePath)

		fileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
		outputFilePathLoudnorm := fileName + ".loudnorm.mkv"
		outputFilePathPeaknorm := fileName + ".peakloud.mkv"
		outputFilePathNormalized := fileName + ".normalized.mkv"

		_, err := normalizer.Loudnorm(inputFilePath, outputFilePathLoudnorm)
		if err != nil {
			return err
		}

		stats, err := normalizer.Analyze(inputFilePath)
		if err != nil {
			return err
		}

		max_volume, err := strconv.ParseFloat(stats.Max_volume[:len(stats.Max_volume)-3], 64)
		if err != nil {
			return err
		}

		err = normalizer.Peaknorm(inputFilePath, outputFilePathPeaknorm, max_volume)
		if err != nil {
			return err
		}

		err = os.Remove(outputFilePathLoudnorm)
		if err != nil {
			return err
		}

		err = os.Rename(outputFilePathPeaknorm, outputFilePathNormalized)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleAnalyze(c *cli.Context) error {
	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		stats, err := normalizer.Analyze(inputFilePath)
		if err != nil {
			return err
		}

		fmt.Printf("       file: %s\n", inputFilePath)
		fmt.Printf("mean_volume: %s\n", stats.Mean_volume)
		fmt.Printf(" max_volume: %s\n\n", stats.Max_volume)
	}
	return nil
}
