package utils

type AllDevice struct {
	UUID         string `json:"device_id"`
	DeviceAddr   string `json:"device_addr"`
	Timestamp    string `json:"timestamp"`
	IsRandomized bool   `json:"is_randomized"`
	DevicePower  int    `json:"device_power"`
	SSID         string `json:"ssid"`
	FrameType    string `json:"frame_type"`
	Zone         string `json:"library_section"`
}
