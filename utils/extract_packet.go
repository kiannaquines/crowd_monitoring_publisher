package utils

import (
	"encoding/json"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"fmt"
	"os"
)

func ExtractPacketInformation(packet gopacket.Packet) {
	if dot11Layer := packet.Layer(layers.LayerTypeDot11); dot11Layer != nil {
		dot11, ok := dot11Layer.(*layers.Dot11)

		if !ok {
			return
		}

		var clientAddr string
		var isRandomized bool
		var signalStrength int
		var ssid string
		var frame string

		timestamp := packet.Metadata().Timestamp.Format("2006-01-02 15:04:05")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		zone := os.Getenv("MQTT_BROKER_ZONE")
		topic := os.Getenv("MQTT_BROKER_TOPIC")

		switch dot11.Type {

		case layers.Dot11TypeDataQOSData:
            if dot11.Flags.ToDS() {
				clientAddr = dot11.Address2.String()
				frame = "QoS Data Frame"
				ssid = "Not Associated"
			} else {
				return
			}
				
		case layers.Dot11TypeMgmtAssociationResp:
			clientAddr = dot11.Address2.String()
			frame = "Probe Response"
			if probeResLayer := packet.Layer(layers.LayerTypeDot11MgmtProbeResp); probeResLayer != nil {
				ssid = extractSSID(probeResLayer.LayerPayload())
			}

		case layers.Dot11TypeMgmtDisassociation:
			clientAddr = dot11.Address1.String()
			frame = "Diassociation"
			if diAssocLayer := packet.Layer(layers.LayerTypeDot11MgmtDisassociation); diAssocLayer != nil {
				ssid = extractSSID(diAssocLayer.LayerPayload())
			}

		case layers.Dot11TypeMgmtAuthentication:
			clientAddr = dot11.Address1.String()
			frame = "Authentication"
			if authLayer := packet.Layer(layers.LayerTypeDot11MgmtAuthentication); authLayer != nil {
				ssid = extractSSID(authLayer.LayerPayload())
			}

		case layers.Dot11TypeMgmtDeauthentication:
			clientAddr = dot11.Address1.String()
			frame = "Deauthentication"
			if deAuthLayer := packet.Layer(layers.LayerTypeDot11MgmtDeauthentication); deAuthLayer != nil {
				ssid = extractSSID(deAuthLayer.LayerPayload())
			}

		case layers.Dot11TypeMgmtAssociationReq:
			clientAddr = dot11.Address2.String()
			frame = "Association Request"
			if assocReqLayer := packet.Layer(layers.LayerTypeDot11MgmtAssociationReq); assocReqLayer != nil {
				ssid = extractSSID(assocReqLayer.LayerPayload())
			}

		case layers.Dot11TypeMgmtProbeReq:
			clientAddr = dot11.Address2.String()
			frame = "Probe Request"
			if probeReqLayer := packet.Layer(layers.LayerTypeDot11MgmtProbeReq); probeReqLayer != nil {
				ssid = extractSSID(probeReqLayer.LayerPayload())
			}

		default:
			return
		}

		isRandomized = isMACRandomized(clientAddr)
		signalStrength, _ = extractSignalStrength(packet)

		device := AllDevice{
			UUID:         uuid.New().String(),
			DeviceAddr:   clientAddr,
			IsRandomized: isRandomized,
			DevicePower:  signalStrength,
			SSID:         ssid,
			Timestamp:    timestamp,
			FrameType:    frame,
			Zone:         zone,
		}

		jsonData, err := json.Marshal(device)

		if err != nil {
			return
		}
		
		token := mqttClient.Publish(topic, 0, false, jsonData)
		token.Wait()

		if token.Error() != nil {
			log.Printf("Failed to publish device data: %v\n", token.Error())
		}
	}

	return
}
