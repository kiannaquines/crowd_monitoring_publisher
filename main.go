package main

import (
	"fmt"
	"log"
	"os"

	"publisher/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/joho/godotenv"
)

func main() {
	utils.MqttInit()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	iface := os.Getenv("MQTT_BROKER_IFACE")

	log.Printf("Crowd Monitoring Publisher is starting at interface %v", iface)

	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		fmt.Println("Starting packet capture indefinitely...")
		capturePackets(packetSource)
	}
}

func capturePackets(packetSource *gopacket.PacketSource) {
	for {
		packet := <-packetSource.Packets()
		utils.ExtractPacketInformation(packet)
	}
}
