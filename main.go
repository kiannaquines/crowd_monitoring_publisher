package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/joho/godotenv"
	"publisher/utils"
)

func main() {
	
	utils.MqttInit()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	iface := os.Getenv("MQTT_BROKER_IFACE")
	captureMin := os.Getenv("MQTT_BROKER_CAPTURE_MINUTE")
	pauseMin := os.Getenv("MQTT_BROKER_PAUSE_MINUTE")

	log.Printf("Crowd Monitoring Publisher is starting at interface %v", iface)

	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	captureDurationMin, err := strconv.Atoi(captureMin)
	if err != nil {
		log.Fatalf("Invalid capture duration: %v", err)
	}
	pauseDurationMin, err := strconv.Atoi(pauseMin)
	if err != nil {
		log.Fatalf("Invalid pause duration: %v", err)
	}

	captureDuration := time.Duration(captureDurationMin) * time.Minute
	pauseDuration := time.Duration(pauseDurationMin) * time.Minute

	for {
		
		fmt.Println("Starting packet capture for 10 minutes...")
		captureUntil := time.After(captureDuration)
		capturePackets(packetSource, captureUntil)

		fmt.Println("Pausing for 5 minutes...")
		pauseUntil := time.After(pauseDuration)
		
		<-pauseUntil
			fmt.Println("Resuming packet capture...")
	}	
}

func capturePackets(packetSource *gopacket.PacketSource, stopSignal <-chan time.Time) {
	for {
		select {
			case packet := <-packetSource.Packets():
				utils.ExtractPacketInformation(packet)

			case <-stopSignal:
				fmt.Println("10 minutes elapsed, stopping capture...")
				return
		}
	}
}

