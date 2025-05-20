package normalizer

import (
	"encoding/json"
	"fmt"
	"normalizer/internal/ffmpeg"
	"os"
	"regexp"
)

type loudnormStats struct {
	Input_i      string
	Input_tp     string
	Input_lra    string
	Input_thresh string

	Output_i      string
	Output_tp     string
	Output_lra    string
	Output_thresh string

	Normalization_type string
	Target_offset      string
}

// first pass of loudnorm, get stats
// ffmpeg -hide_banner -i $INPUT -c:v copy -c:a copy -b:a 256k -filter:a loudnorm=print_format=json -f null /dev/null

// second pass of loudnorm, apply stats
//
//	ffmpeg -hide_banner -i $INPUT -c:v copy -c:a copy -b:a 256k \
//	    -filter:a loudnorm=linear=true:measured_I=$INPUT_I:measured_LRA=$INPUT_LRA:measured_tp=$INPUT_TP:measured_thresh=$INPUT_THRESH \
//	    $OUTPUT_NORMALIZED
func Loudnorm(filePath string, outputFilePath string) (*loudnormStats, error) {
	output, err := ffmpeg.Run(
		ffmpeg.WithInput(filePath),
		ffmpeg.WithVideoCodec("copy"),
		ffmpeg.WithAudioFilter("loudnorm=print_format=json"),
		ffmpeg.WithOverwrite(),
		ffmpeg.WithFormat("null"),
		ffmpeg.WithOutput(os.DevNull),
	)

	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	stats, err := parseLoudnormStats(string(output))
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to parse loudnorm stats: %w", err)
	}

	audioFilter := fmt.Sprintf(
		"loudnorm=linear=true:measured_I=%v:measured_LRA=%v:measured_tp=%v:measured_thresh=%v:print_format=json",
		stats.Input_i, stats.Input_lra, stats.Input_tp, stats.Input_thresh)

	output, err = ffmpeg.Run(
		ffmpeg.WithInput(filePath),
		ffmpeg.WithVideoCodec("copy"),
		ffmpeg.WithAudioFilter(audioFilter),
		ffmpeg.WithOverwrite(),
		ffmpeg.WithOutput(outputFilePath),
	)

	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	outputStats, err := parseLoudnormStats(string(output))
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to parse loudnorm stats: %w", err)
	}

	return outputStats, nil
}

var jsonRegex = regexp.MustCompile(`{([\w\n\s\:,\-."]*)}`)

func parseLoudnormStats(output string) (*loudnormStats, error) {
	jsonPart := jsonRegex.Find([]byte(output))

	var stats loudnormStats
	err := json.Unmarshal(jsonPart, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &stats, nil
}
