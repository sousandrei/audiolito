package normalizer

import (
	"fmt"
	"normalizer/internal/ffmpeg"
	"os"
	"regexp"
	"strconv"
)

type volumeStats struct {
	Mean_volume float64
	Max_volume  float64
}

var dbfloatRegex = regexp.MustCompile(`: -([0-9]{1,2}.[0-9]{1,2})`)

// ffmpeg -hide_banner -i $OUTPUT_NORMALIZED -filter:a volumedetect -f matroska /dev/null
func Analyze(filePath string) (*volumeStats, error) {
	output, err := ffmpeg.Run(
		ffmpeg.WithInput(filePath),
		ffmpeg.WithVideoCodec("copy"),
		ffmpeg.WithAudioFilter("volumedetect"),
		ffmpeg.WithOverwrite(),
		ffmpeg.WithFormat("matroska"),
		ffmpeg.WithOutput(os.DevNull),
	)

	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	matches := dbfloatRegex.FindAllStringSubmatch(string(output), 2)

	meanVolume, err := strconv.ParseFloat(matches[0][1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mean volume: %w", err)
	}

	maxVolume, err := strconv.ParseFloat(matches[1][1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse max volume: %w", err)
	}

	return &volumeStats{
		Mean_volume: meanVolume,
		Max_volume:  maxVolume,
	}, nil
}
