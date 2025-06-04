package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type ffmpeg struct {
	args  []string
	cmd   string
	debug bool
}

func Run(ctx context.Context, options ...func(*ffmpeg)) ([]byte, error) {
	f := &ffmpeg{
		cmd: "ffmpeg",
		args: []string{
			"-hide_banner",
		},
	}

	for _, option := range options {
		option(f)
	}

	cmd := exec.CommandContext(ctx, f.cmd, f.args...)

	output, err := runCommand(cmd, f.debug)
	if err != nil {
		return nil, fmt.Errorf("failed to run ffmpeg command: %w", err)
	}

	return output, nil
}

func runCommand(cmd *exec.Cmd, debug bool) ([]byte, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	var combinedOutput []byte

	var wg sync.WaitGroup

	for _, stream := range []struct {
		pipe       io.Reader
		streamName string
	}{
		{stdout, "stdout"},
		{stderr, "stderr"},
	} {
		wg.Add(1)
		go func(pipe io.Reader, streamName string) {
			defer wg.Done()

			scanner := bufio.NewScanner(pipe)
			for scanner.Scan() {
				combinedOutput = append(combinedOutput, scanner.Bytes()...)

				if debug {
					fmt.Println(scanner.Text())
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading from %s pipe: %v\n", streamName, err)
			}
		}(stream.pipe, stream.streamName)
	}

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	wg.Wait()

	return combinedOutput, nil
}

func WithCommand(command string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.cmd = command
	}
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
		if output == "" {
			return
		}
		d.args = append(d.args, output)
	}
}

func WithOverwrite() func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.args = append(d.args, "-y")
	}
}

func WithProgress(addr string) func(*ffmpeg) {
	return func(d *ffmpeg) {
		if addr == "" {
			return
		}
		d.args = append(d.args, "-progress", addr)
	}
}

func WithDebug(debug bool) func(*ffmpeg) {
	return func(d *ffmpeg) {
		d.debug = debug
	}
}
