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

		err = normalizer.Peaknorm(inputFilePath, outputFilePathPeaknorm, stats.Max_volume)
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
