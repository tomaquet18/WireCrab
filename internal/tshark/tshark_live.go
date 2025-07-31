package tshark

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type TsharkLive struct {
	cmd         *exec.Cmd
	stdout      *bufio.Reader
	stderr      io.ReadCloser
	pcapPath    string
	mu          sync.Mutex
	closed      bool
	packetsLock sync.RWMutex
	packets     []ProtocolInfo
}

func StartTsharkLive(device string) (*TsharkLive, error) {
	binary := "tshark"
	if runtime.GOOS == "windows" {
		binary = "tshark.exe"
	}

	// Create a temporary file for pcap
	tempDir := os.TempDir()
	pcapPath := filepath.Join(tempDir, fmt.Sprintf("wirecrab-%d.pcap", time.Now().Unix()))

	cmd := exec.Command(binary,
		"-i", device,
		"-l",           // line-buffered
		"-n",           // no name resolution
		"-w", pcapPath, // write to pcap file
		"-T", "fields", // fields output
		"-E", "separator=|",
		"-e", "frame.number",
		"-e", "_ws.col.protocol",
		"-e", "ip.src",
		"-e", "ip.dst",
		"-e", "frame.len",
		"-e", "_ws.col.info")

	// Hide window on Windows
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	t := &TsharkLive{
		cmd:      cmd,
		stdout:   bufio.NewReader(stdout),
		stderr:   stderr,
		pcapPath: pcapPath,
		packets:  make([]ProtocolInfo, 0),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start tshark: %w", err)
	}

	// Start background packet collector
	go t.collectPackets()

	return t, nil
}

func (t *TsharkLive) collectPackets() {
	for {
		proto, err := t.Next()
		if err != nil {
			if err == io.EOF {
				log.Printf("Error reading packet: %v", err)
			}
			return
		}

		t.packetsLock.Lock()
		t.packets = append(t.packets, *proto)
		t.packetsLock.Unlock()
	}
}

func (t *TsharkLive) Next() (*ProtocolInfo, error) {
	line, err := t.stdout.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("read stdout: %w", err)
	}

	// Split the line using pipe separator
	fields := strings.Split(strings.TrimSpace(line), "|")
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid number of fields")
	}

	// Create a ProtocolInfo with the fields
	return &ProtocolInfo{
		Name: fields[1], // protocol
		Detail: map[string]any{
			"frame.number": map[string]any{"value": fields[0]},
			"ip.src":       map[string]any{"value": fields[2]},
			"ip.dst":       map[string]any{"value": fields[3]},
			"frame.len":    map[string]any{"value": fields[4]},
			"_ws.col.info": map[string]any{"value": fields[5]},
		},
		Child: nil,
	}, nil
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
	// Clean up pcap file
	_ = os.Remove(t.pcapPath)
	return t.cmd.Wait()
}

func (t *TsharkLive) GetPacketList(offset, limit int) ([]ProtocolInfo, error) {
	t.packetsLock.RLock()
	defer t.packetsLock.RUnlock()

	if offset >= len(t.packets) {
		return nil, nil
	}

	end := offset + limit
	if end > len(t.packets) {
		end = len(t.packets)
	}

	return t.packets[offset:end], nil
}

func (t *TsharkLive) GetPacketCount() (int, error) {
	t.packetsLock.RLock()
	defer t.packetsLock.RUnlock()
	return len(t.packets), nil
}

// GetPacketDetails retrieves detailed information for a specific packet
func (t *TsharkLive) GetPacketDetails(packetNumber int) (*PacketDetails, error) {
	binary := "tshark"
	if runtime.GOOS == "windows" {
		binary = "tshark.exe"
	}

	// Hide window on Windows
	sysAttrs := &syscall.SysProcAttr{
		HideWindow: true,
	}

	// Get PDML data
	pdmlCmd := exec.Command(binary,
		"-r", t.pcapPath, // read from pcap file
		"-Y", fmt.Sprintf("frame.number==%d", packetNumber), // filter by packet number
		"-T", "pdml") // use PDML for detailed output
	pdmlCmd.SysProcAttr = sysAttrs

	pdmlOutput, err := pdmlCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get packet details: %w", err)
	}

	// Get hex dump
	hexCmd := exec.Command(binary,
		"-r", t.pcapPath,
		"-Y", fmt.Sprintf("frame.number==%d", packetNumber),
		"-x")
	hexCmd.SysProcAttr = sysAttrs

	hexOutput, err := hexCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get hex dump: %w", err)
	}

	// Parse PDML output and return ProtocolInfo
	var pdml PDML
	if err := xml.Unmarshal(pdmlOutput, &pdml); err != nil {
		return nil, fmt.Errorf("parse pdml: %w", err)
	}

	if len(pdml.Packets) == 0 {
		return nil, fmt.Errorf("packet not found")
	}

	return &PacketDetails{
		Info:    PdmlToProtocolInfo(pdml.Packets[0].Protos),
		HexDump: string(hexOutput),
	}, nil
}
