package capture

import (
	"wirecrab/internal/types"
)

type CaptureSession struct {
	Packets []types.CapturedPacket
}
