package ffmpeg

import (
	"fmt"
	"os/exec"
)

var defaultArgs = []string{
	"-hide_banner",
}

type ffmpeg struct {
	args []string
}

func Run(options ...func(*ffmpeg)) ([]byte, error) {
	d := &ffmpeg{
		args: defaultArgs,
	}

	for _, option := range options {
		option(d)
	}

	cmd := exec.Command("ffmpeg", d.args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return out, nil
}

func WithInput(input string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-i", input)
	}
}

func WithVideoCodec(codec string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-c:v", codec)
	}
}

func WithAudioCodec(codec string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-c:a", codec)
	}
}

func WithAudioBitrate(bitrate string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-b:a", bitrate)
	}
}

func WithAudioFilter(filter string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-filter:a", filter)
	}
}

func WithFormat(format string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-f", format)
	}
}

func WithOutput(output string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, output)
	}
}

func WithOverwrite() func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-y")
	}
}
