package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/joho/godotenv"
)

var blockedOUIs map[string]struct{}

func LoadBlockedOUIs(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening OUI file: %v", err)
	}
	defer file.Close()

	blockedOUIs = make(map[string]struct{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		oui := strings.ToLower(scanner.Text())
		blockedOUIs[oui] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading OUI file: %v", err)
	}
}

func isOUIBlocked(macAddr string) bool {
	oui := macAddr[:8]
	_, blocked := blockedOUIs[oui]
	return blocked
}

func ExtractPacketInformation(packet gopacket.Packet) {
	if dot11Layer := packet.Layer(layers.LayerTypeDot11); dot11Layer != nil {
		dot11, ok := dot11Layer.(*layers.Dot11)

		if !ok {
			return
		}

		var clientAddr string
		var isRandomized bool
		var signalStrength int
		var frame string

		timestamp := packet.Metadata().Timestamp.Format("2006-01-02 15:04:05")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		zone := os.Getenv("MQTT_BROKER_ZONE")
		topic := os.Getenv("MQTT_BROKER_TOPIC")

		LoadBlockedOUIs("oui.txt")

		switch dot11.Type {
		case layers.Dot11TypeDataQOSData:
			if dot11.Flags.ToDS() && uint8(dot11.Flags) == 0x41 {
				clientAddr = dot11.Address2.String()
				frame = "QoS Data Frame"
			} else {
				return
			}

		case layers.Dot11TypeMgmtAuthentication:
			if dot11.Flags.ToDS() && uint8(dot11.Flags) == 0x41 {
				clientAddr = dot11.Address1.String()
				frame = "Authentication"
			}

		case layers.Dot11TypeMgmtAssociationReq:
			if dot11.Flags.ToDS() && uint8(dot11.Flags) == 0x41 {
				clientAddr = dot11.Address2.String()
				frame = "Association Request"
			}

		case layers.Dot11TypeMgmtProbeReq:
			clientAddr = dot11.Address2.String()
			frame = "Probe Request"

		default:
			return
		}

		if isOUIBlocked(clientAddr) {
			fmt.Printf("Blocked %s \n", clientAddr)
			return
		}

		isRandomized = isMACRandomized(clientAddr)
		signalStrength, isValid := extractSignalStrength(packet)
		if isValid {
			if signalStrength >= -80 && signalStrength <= -30 {
				device := AllDevice{
					DeviceAddr:   clientAddr,
					IsRandomized: isRandomized,
					DevicePower:  signalStrength,
					Timestamp:    timestamp,
					FrameType:    frame,
					Zone:         zone,
				}

				jsonData, err := json.Marshal(device)

				if err != nil {
					return
				}

				fmt.Printf("%s %d %s \n", clientAddr, signalStrength, frame)
				token := mqttClient.Publish(topic, 0, false, jsonData)
				token.Wait()

				if token.Error() != nil {
					log.Printf("Failed to publish device data: %v\n", token.Error())
				}
			}
		}
	}
}
