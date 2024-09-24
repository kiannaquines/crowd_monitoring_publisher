package utils

func extractSSID(ie []byte) string {
	if len(ie) > 2 && ie[0] == 0x00 {
		ssidLength := int(ie[1])
		ssid := ie[2 : 2+ssidLength]

		if len(ssid) > 0 {
			return string(ssid)
		} else {
			return "Hidden"
		}
	}
	return "Not Associated"
}
