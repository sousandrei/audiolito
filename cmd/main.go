package main

import (
	"fmt"
	"log"
	"normalizer/internal/ffmpeg"
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
			{
				Name:      "wav",
				ArgsUsage: "[FILES]",
				Action:    handleWav,
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

		fmt.Printf("loud normalizing %s\n", inputFilePath)

		fileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
		extension := filepath.Ext(inputFilePath)
		outputFilePathLoudNormalized := fileName + ".loud_normalized" + extension
		outputFilePathNormalized := fileName + ".normalized" + extension

		stats, err := normalizer.Loudnorm(inputFilePath, outputFilePathLoudNormalized)
		if err != nil {
			return err
		}

		fmt.Printf("peak normalizing %s\n", outputFilePathLoudNormalized)

		inputI, err := strconv.ParseFloat(stats.Input_i, 64)
		if err != nil {
			return fmt.Errorf("failed to parse input_i: %w", err)
		}

		targetLoudness := inputI * -1

		err = normalizer.Peaknorm(outputFilePathLoudNormalized, outputFilePathNormalized, targetLoudness)
		if err != nil {
			return fmt.Errorf("failed to apply peak normalization: %w", err)
		}

		if err := os.Remove(outputFilePathLoudNormalized); err != nil {
			return fmt.Errorf("failed to remove temporary file %s: %w", outputFilePathLoudNormalized, err)
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
		fmt.Printf("mean_volume: %s\n", strconv.FormatFloat(stats.Mean_volume, 'f', -1, 64))
		fmt.Printf(" max_volume: %s\n\n", strconv.FormatFloat(stats.Max_volume, 'f', -1, 64))
	}
	return nil
}

func handleWav(c *cli.Context) error {
	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		fmt.Printf("converting %s to wav\n", inputFilePath)

		fileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
		outputFilePath := fileName + ".wav"

		_, err := ffmpeg.Run(
			ffmpeg.WithInput(inputFilePath),
			ffmpeg.WithAudioCodec("pcm_s16le"),
			ffmpeg.WithOverwrite(),
			ffmpeg.WithOutput(outputFilePath),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
