package normalizer

import (
	"context"
	"fmt"
	"normalizer/internal/ffmpeg"
	"normalizer/internal/tui"
	"os"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type volumeStats struct {
	Mean_volume float64
	Max_volume  float64
}

var dbfloatRegex = regexp.MustCompile(`: (-?[0-9]{1,2}.[0-9]{1,2}) dB`)

// ffmpeg -hide_banner -i $OUTPUT_NORMALIZED -filter:a volumedetect -f null /dev/null
func Analyze(
	ctx context.Context,
	options ...func(*normalizer),
) (*volumeStats, error) {
	n := &normalizer{}

	for _, option := range options {
		option(n)
	}

	if n.filePath == "" {
		return nil, fmt.Errorf("file path is required")
	}

	stats, err := getStats(ctx, n.filePath, n.debug)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	var p *tea.Program
	var addr string

	if n.debug {
		n.tui = false // disable TUI in debug mode
	}

	if n.tui {
		pn, addrn, err := tui.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create TUI: %w", err)
		}
		p = pn
		addr = addrn

		go p.Run()

		p.Send(tui.SetDurationMsg(stats.duration))
	} else {
		addr = ""
	}

	output, err := ffmpeg.Run(
		ctx,
		ffmpeg.WithInput(n.filePath),
		ffmpeg.WithVideoCodec("copy"),
		ffmpeg.WithAudioFilter("volumedetect"),
		ffmpeg.WithOverwrite(),
		ffmpeg.WithFormat("null"),
		ffmpeg.WithOutput(os.DevNull),
		ffmpeg.WithProgress(addr),
		ffmpeg.WithDebug(n.debug),
	)
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	matches := dbfloatRegex.FindAllStringSubmatch(string(output), 2)

	if len(matches) != 2 {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to parse volume stats")
	}

	meanVolume, err := strconv.ParseFloat(matches[0][1], 64)
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to parse mean volume: %w", err)
	}

	maxVolume, err := strconv.ParseFloat(matches[1][1], 64)
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("failed to parse max volume: %w", err)
	}

	return &volumeStats{
		Mean_volume: meanVolume,
		Max_volume:  maxVolume,
	}, nil
}
