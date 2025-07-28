package tshark

import (
	"encoding/xml"
	"strings"
)

type PDML struct {
	XMLName xml.Name     `xml:"pdml"`
	Packets []PDMLPacket `xml:"packet"`
}

type PDMLPacket struct {
	Protos []PDMLProto `xml:"proto"`
}

type PDMLProto struct {
	Name   string      `xml:"name,attr"`
	Fields []PDMLField `xml:"field"`
	Protos []PDMLProto `xml:"proto"`
}

type PDMLField struct {
	Name string `xml:"name,attr"`
	Show string `xml:"show,attr"`
	Pos  string `xml:"pos,attr"`
	Size string `xml:"size,attr"`
}

type ProtocolInfo struct {
	Name   string
	Detail interface{}
	Child  *ProtocolInfo
}

func PdmlToProtocolInfo(protos []PDMLProto) *ProtocolInfo {
	if len(protos) == 0 {
		return nil
	}
	return &ProtocolInfo{
		Name:   protos[0].Name,
		Detail: pdmlFieldsToMap(protos[0].Fields),
		Child:  PdmlToProtocolInfo(protos[1:]),
	}
}

func pdmlFieldsToMap(fields []PDMLField) map[string]any {
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
