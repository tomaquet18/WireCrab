package tshark

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

// TsharkLive wraps a long-running tshark process that captures and dissects.
type TsharkLive struct {
	cmd    *exec.Cmd
	stdout *bufio.Reader
	stderr io.ReadCloser

	mu     sync.Mutex
	closed bool
}

// StartTsharkLive starts tshark capturing on the given device and emitting PDML continiously.
func StartTsharkLive(device string) (*TsharkLive, error) {
	binary := "tshark"
	if runtime.GOOS == "windows" {
		binary = "tshark.exe"
	}

	// -l : line-buffered -> we can read incrementally
	// -n : no name resolution (faster)
	// -T pdml : structured XML output
	cmd := exec.Command(binary, "-i", device, "-l", "-n", "-T", "pdml")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start tshark: %w", err)
	}

	return &TsharkLive{
		cmd:    cmd,
		stdout: bufio.NewReader(stdout),
		stderr: stderr,
	}, nil
}

// Next reads until it finds the next </packet> and returns its parsed ProtocolInfo.
// Returns (nil, io.EOF) when tshark finishes.
func (t *TsharkLive) Next() (*ProtocolInfo, error) {
	var buf bytes.Buffer

	// We want to read PDML packet by packet.
	// tshark emits <pdml> ... </packet> ... </packet> ... </pdml> chunks.
	// We accomulate until we hit </packet> or </pdml>
	for {
		line, err := t.stdout.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("read stdout: %w", err)
		}
		buf.Write(line)

		// stop when we got a closing tag of packet
		if bytes.Contains(line, []byte("</packet>")) {
			break
		}

		// sometimes tshark might flush </pdml> too
		if bytes.Contains(line, []byte("</pdml>")) && buf.Len() > 0 {
			break
		}
	}

	// Try to unmarshal the buffer. It can be either a full <pdml>...</pdml> or at least a <packet>...</packet>.
	// Wrap it in <pdml> if we only got <packet>...</packet>.
	data := buf.Bytes()
	if !bytes.Contains(data, []byte("<pdml")) {
		data = append([]byte("<pdml>"), data...)
		data = append(data, []byte("</pdml>")...)
	}

	var result PDML
	if err := xml.Unmarshal(data, &result); err != nil {
		// Log and skip malformed chunks instead of crashing everything
		return nil, fmt.Errorf("xml unmarshall: %w", err)
	}
	if len(result.Packets) == 0 {
		return nil, fmt.Errorf("no packets found in chunk")
	}

	return PdmlToProtocolInfo(result.Packets[0].Protos), nil
}

// Close kills tshark.
func (t *TsharkLive) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	if t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
	}
	return t.cmd.Wait()
}
