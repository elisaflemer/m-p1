package main

import (
	"fmt"
	"bytes"
	"flag"
	"os"
	"encoding/json"
	"net/http"
	"time"
	//"github.com/go-chi/chi"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConnector interface {
	Connect(nodeName string, username string, password string) MQTT.Client
}

type Configuration struct {
	Unit             string  `json:"unit"`
	TransmissionRate float64 `json:"transmission_rate_hz"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	Sensor           string  `json:"sensor"`
	QoS			     byte    `json:"qos"`
}

type Data struct {
	Value            float64   `json:"value"`
	Unit             string    `json:"unit"`
	TransmissionRate float64   `json:"transmission_rate"`
	Longitude        float64   `json:"longitude"`
	Latitude         float64   `json:"latitude"`
	Sensor           string    `json:"sensor"`
	Timestamp        time.Time `json:"timestamp"`
	QoS			  byte      `json:"qos"`
}

func readConfig(filename string) (Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	return config, err
}


type LocalMQTTConnector struct{}

func (l *LocalMQTTConnector) Connect(nodeName string, username string, password string) MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID(nodeName)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

type HiveMQConnector struct{}

func (h *HiveMQConnector) Connect(nodeName string, username string, password string) MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker("tls://b9f3c31144f64d469f184727678d8fb6.s1.eu.hivemq.cloud:8883/mqtt")
	opts.SetClientID(nodeName)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := MQTT.NewClient(opts)
	fmt.Println("Connecting to HiveMQ")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connected to HiveMQ")
	return client
}


func postStructAsJSON(url string, data []byte) error {

	buffer := bytes.NewBuffer(data)
	// Make a POST request
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Println("POST request successful")
	return nil
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("ola")
	fmt.Printf("Recebido: %s do tópico: %s\n", msg.Payload(), msg.Topic())

	postStructAsJSON("http://localhost:5000/api", msg.Payload())

}

func main() {

	configPath := flag.String("config", "", "Path to the configuration file")
	connection := flag.String("connection", "hivemq", "Enter 'hivemq' or 'local' for MQTT connection")

	hivemqUsername := flag.String("username", "", "HiveMQ username")
	hivemqPassword := flag.String("password", "", "HiveMQ password")

	flag.Parse()


	if *configPath == "" {
		fmt.Println("Usage: go run subscriber.go -config <config_path> -connection <hivemq/local> -username <username> -password <password>")
		return
	}

	var connector MQTTConnector

	if *connection == "hivemq"{
		connector = &HiveMQConnector{}
	} else if *connection == "local" {
		connector = &LocalMQTTConnector{}
	} else {
		fmt.Println("Invalid connection type. Enter 'hivemq' or 'local'")
		return	
	}


	client := connector.Connect("subscriber", *hivemqUsername, *hivemqPassword)
	defer client.Disconnect(250)

	if token := client.Subscribe("sensor/solar", 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Println("Subscriber está rodando. Pressione CTRL+C para sair.")
	select {} // Bloqueia indefinidamente
}