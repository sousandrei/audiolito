package normalizer

import (
	"fmt"
	"normalizer/internal/ffmpeg"
	"os"
	"regexp"
)

type volumeStats struct {
	Mean_volume string
	Max_volume  string
}

var meanVolumeRegex = regexp.MustCompile(`mean_volume: (.+)`)
var maxVolumeRegex = regexp.MustCompile(`max_volume: (.+)`)

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

	meanVolume := meanVolumeRegex.FindStringSubmatch(string(output))[1]
	maxVolume := maxVolumeRegex.FindStringSubmatch(string(output))[1]

	return &volumeStats{
		Mean_volume: meanVolume,
		Max_volume:  maxVolume,
	}, nil
}
