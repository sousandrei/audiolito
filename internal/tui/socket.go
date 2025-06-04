package tui

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type progressServer struct {
	addr string
	p    *tea.Program
}

func newServer(ctx context.Context, p *tea.Program) (*progressServer, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket: %w", err)
	}

	// ensure the listener is closed when the context is done
	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close listener: %v\n", err)
		}
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// If the error is due to listener being closed, exit the loop
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Fprintf(os.Stderr, "failed to accept connection: %v\n", err)
					return
				}
			}

			go handleRequest(p, conn)
		}
	}()

	addr := fmt.Sprintf("tcp://%s", listener.Addr().String())

	ps := &progressServer{addr, p}

	return ps, nil
}

func handleRequest(p *tea.Program, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		lines := string(buf[:n])

		updates := map[string]string{}

		for _, line := range strings.Split(lines, "\n") {
			if len(line) == 0 {
				continue
			}

			parts := strings.Split(line, "=")
			updates[parts[0]] = parts[1]
		}

		for _, msg := range processUpdates(updates) {
			p.Send(msg)
		}
	}
}

func processUpdates(updates map[string]string) []tea.Msg {
	// bitrate:N/A
	// drop_frames:0
	// dup_frames:0
	// out_time:00:00:05.943175
	// out_time_ms:5943175
	// out_time_us:5943175
	// progress:end
	// speed:3.15e+03x
	// total_size:N/A

	// bitrate:N/A
	// drop_frames:0
	// dup_frames:0
	// fps:159244.42
	// frame:562784
	// out_time:02:36:19.690667
	// out_time_ms:9379690667
	// out_time_us:9379690667
	// progress:continue
	// speed:2.65e+03x
	// stream_0_0_q:-1.0
	// total_size:N/A

	var msgs []tea.Msg

	for key, value := range updates {
		switch key {
		case "progress":
			if value == "end" {
				msgs = append(msgs, tea.QuitMsg{})
				continue
			}

		case "out_time":
			var hours, minutes, seconds float64
			n, err := fmt.Sscanf(value, "%f:%f:%f", &hours, &minutes, &seconds)
			if err != nil {
				fmt.Printf("Error parsing duration: %v\n", err)
				continue
			}
			if n != 3 {
				fmt.Printf("Expected 3 parsed items, got %d. Invalid format.\n", n)
				continue
			}

			outTime := hours*3600 + minutes*60 + seconds

			msgs = append(msgs, ProgressMsg(outTime))
		}

	}

	return msgs
}
