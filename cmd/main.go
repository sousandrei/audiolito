package main

import (
	"fmt"
	"log"
	"normalizer/internal/ffmpeg"
	"normalizer/internal/normalizer"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "audiolito",
		Usage: "tools for audio manipulation",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-tui",
				Usage: "disables the TUI",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "pipes ffmpeg output to stdout",
				Value: false,
			},
		},
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
	ctx := c.Context

	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		fmt.Printf("normalizing %s\n", inputFilePath)

		fileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
		extension := filepath.Ext(inputFilePath)
		outputFilePathNormalized := fileName + ".normalized" + extension

		_, err := normalizer.Loudnorm(
			ctx,
			inputFilePath,
			outputFilePathNormalized,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleWav(c *cli.Context) error {
	ctx := c.Context

	for _, inputFilePath := range c.Args().Slice() {
		if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", inputFilePath)
		}

		fmt.Printf("converting %s to wav\n", inputFilePath)

		fileName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
		outputFilePath := fileName + ".wav"

		_, err := ffmpeg.Run(
			ctx,
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
