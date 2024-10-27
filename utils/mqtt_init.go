package utils

import (
	"crypto/tls"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

var mqttClient mqtt.Client

func MqttInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("MQTT_BROKER_USERNAME")
	password := os.Getenv("MQTT_BROKER_PASSWORD")
	clientID := os.Getenv("MQTT_BROKER_CLIENTID")
	protocol := os.Getenv("MQTT_BROKER_PROTOCOL")
	host := os.Getenv("MQTT_BROKER_HOST")
	port := os.Getenv("MQTT_BROKER_PORT")

	uri := protocol + host + ":" + port
	opts := mqtt.NewClientOptions().AddBroker(uri)

	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	opts.SetTLSConfig(tlsConfig)

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v\n", err)
	}

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}
