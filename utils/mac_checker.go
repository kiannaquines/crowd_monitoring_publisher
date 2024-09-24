package utils

import (
	"log"
	"strings"
)

func isMACRandomized(macAddress string) bool {
	macAddress = strings.ReplaceAll(macAddress, ":", "")

	if len(macAddress) != 12 {
		log.Println("Invalid MAC address format")
		return false
	}

	firstByte, err := hexToByte(macAddress[:2])
	if err != nil {
		log.Println("Error parsing MAC address:", err)
		return false
	}

	return (firstByte & 0x02) != 0
}
