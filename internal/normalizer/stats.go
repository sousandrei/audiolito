package normalizer

import (
	"context"
	"fmt"
	"normalizer/internal/ffmpeg"
	"regexp"
)

type stats struct {
	duration float64 //duration in ms
}

var durationRegex = regexp.MustCompile(`Duration: (\d{2}:\d{2}:\d{2}\.\d{1,3}),`)

func getStats(
	ctx context.Context,
	filePath string,
	debug bool,
) (*stats, error) {
	output, err := ffmpeg.Run(
		ctx,
		ffmpeg.WithCommand("ffprobe"),
		ffmpeg.WithInput(filePath),
		ffmpeg.WithDebug(debug),
	)
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	matches := durationRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse duration from ffprobe output")
	}

	durationParts := matches[1]
	var hours, minutes, seconds float64
	fmt.Sscanf(durationParts, "%02f:%02f:%02f", &hours, &minutes, &seconds)

	duration := (hours*3600 + minutes*60 + seconds)

	return &stats{
		duration: duration,
	}, nil
}
