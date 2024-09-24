package utils

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func extractSignalStrength(packet gopacket.Packet) (int, bool) {
	radiotapLayer := packet.Layer(layers.LayerTypeRadioTap)
	if radiotapLayer == nil {
		return 0, false
	}

	radiotap, ok := radiotapLayer.(*layers.RadioTap)
	if !ok {
		return 0, false
	}

	return int(radiotap.DBMAntennaSignal), true
}
