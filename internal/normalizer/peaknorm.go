package normalizer

import (
	"fmt"
	"normalizer/internal/ffmpeg"
)

// ffmpeg -hide_banner -i $OUTPUT_NORMALIZED -c:v copy -af "volume=$MAX_VOLUME" $OUTPUT_NORMALIZED_LOUD
func Peaknorm(filePath string, outputFilePath string, targetLevel float64) error {
	audioFilter := fmt.Sprintf("volume=%f", targetLevel)

	output, err := ffmpeg.Run(
		ffmpeg.WithInput(filePath),
		ffmpeg.WithVideoCodec("copy"),
		ffmpeg.WithAudioFilter(audioFilter),
		ffmpeg.WithOverwrite(),
		ffmpeg.WithOutput(outputFilePath),
	)

	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return nil
}
