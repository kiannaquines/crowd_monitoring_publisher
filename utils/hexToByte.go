package utils

import "fmt"

func hexToByte(hexStr string) (byte, error) {
	var byteValue byte
	_, err := fmt.Sscanf(hexStr, "%2X", &byteValue)
	if err != nil {
		return 0, err
	}
	return byteValue, nil
}
