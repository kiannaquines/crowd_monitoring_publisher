package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/joho/godotenv"
	"log"
	"os"
	"publisher/utils"
)

func main() {
	utils.MqttInit()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	iface := os.Getenv("MQTT_BROKER_IFACE")

	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	log.Printf("Crowd Monitoring Publisher is starting at interface %v", iface)
	for packet := range packetSource.Packets() {
		utils.ExtractPacketInformation(packet)
	}
}
