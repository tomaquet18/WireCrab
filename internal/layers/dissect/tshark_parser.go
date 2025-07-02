package dissect

import (
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type pdml struct {
	XMLName xml.Name     `xml:"pdml"`
	Packets []pdmlPacket `xml:"packet"`
}

type pdmlPacket struct {
	Protos []pdmlProto `xml:"proto"`
}

type pdmlProto struct {
	Name   string      `xml:"name,attr"`
	Fields []pdmlField `xml:"field"`
	Protos []pdmlProto `xml:"proto"` // nested protocols
}

type pdmlField struct {
	Name string `xml:"name,attr"`
	Show string `xml:"show,attr"`
	Pos  string `xml:"pos,attr"`
	Size string `xml:"size,attr"`
}

func ParseWithTshark(payload []byte) (*ProtocolInfo, error) {
	tsharkBinary := "tshark"
	if runtime.GOOS == "windows" {
		tsharkBinary = "tshark.exe"
	}
	cmd := exec.Command(tsharkBinary, "-i", "-", "-T", "pdml")
	fmt.Print(cmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err = writePcapTo(stdin, payload)
	stdin.Close()
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	cmd.Wait()

	var result pdml
	if err := xml.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	if len(result.Packets) == 0 {
		return nil, fmt.Errorf("no packets in pdml output")
	}

	return pdmlToProtocolInfo(result.Packets[0].Protos), nil
}

func pdmlToProtocolInfo(protos []pdmlProto) *ProtocolInfo {
	if len(protos) == 0 {
		return nil
	}
	root := &ProtocolInfo{
		Name:   protos[0].Name,
		Detail: pdmlFieldsToMap(protos[0].Fields),
		Child:  pdmlToProtocolInfo(protos[1:]),
	}
	return root
}

func pdmlFieldsToMap(fields []pdmlField) map[string]any {
	m := make(map[string]any)
	for _, f := range fields {
		m[strings.ToLower(f.Name)] = map[string]any{
			"value": f.Show,
			"pos":   f.Pos,
			"size":  f.Size,
		}
	}
	return m
}

// WriteSinglePacketAsPcapTo writes a single packet to an io.Writes in pcap format
func writePcapTo(w io.Writer, data []byte) error {
	writer := pcapgo.NewWriter(w)
	if err := writer.WriteFileHeader(65536, layers.LinkTypeEthernet); err != nil {
		return err
	}

	captureInfo := gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: len(data),
		Length:        len(data),
	}
	return writer.WritePacket(captureInfo, data)
}
