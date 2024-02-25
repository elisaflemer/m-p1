package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
	"flag"
	"github.com/joho/godotenv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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

type MQTTConnector interface {
	Connect(nodeName string) MQTT.Client
}

type LocalMQTTConnector struct{}

func (l *LocalMQTTConnector) Connect(nodeName string) MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID(nodeName)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

type HiveMQConnector struct{}

func (h *HiveMQConnector) Connect(nodeName string) MQTT.Client {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	opts := MQTT.NewClientOptions().AddBroker("tls://b9f3c31144f64d469f184727678d8fb6.s1.eu.hivemq.cloud:8883/mqtt")
	opts.SetClientID(nodeName)
	username := os.Getenv("HIVEMQ_USERNAME")
	password := os.Getenv("HIVEMQ_PASSWORD")
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func main() {
	configPath := flag.String("config", "", "Path to the configuration file")
	csvPath := flag.String("csv", "", "Path to the CSV file")
	connection := flag.String("connection", "hivemq", "Enter 'hivemq' or 'local' for MQTT connection")

	flag.Parse()

	if *configPath == "" || *csvPath == "" {
		fmt.Println("Usage: go run publisher.go -config <config_path> -csv <csv_path> -connection <hivemq/local>")
		return
	}

	config, err := readConfig(*configPath)
	if err != nil {
		panic(err)
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

	client := connector.Connect("publisher")
	defer client.Disconnect(250)

	data, err := readCSV(*csvPath)
	if err != nil {
		panic(err)
	}

	publishData(client, config, data)
}

func connectMQTT(node_name string) MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker("tls://b9f3c31144f64d469f184727678d8fb6.s1.eu.hivemq.cloud:8883/mqtt")
	opts.SetClientID(node_name)
	opts.SetPassword("Levinson1")
	opts.SetUsername("elisa")
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

func readCSV(csvPath string) ([]float64, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []float64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		var value float64
		_, err := fmt.Sscanf(line, "%f", &value)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

func publishData(client MQTT.Client, config Configuration, data []float64) {
	interval := time.Second / time.Duration(config.TransmissionRate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for _, value := range data {
		roundedValue := math.Round(value*100) / 100

		message := createJSONMessage(config, roundedValue)

		token := client.Publish("sensor/"+config.Sensor, byte(config.QoS), false, message)
		//fmt.Fprintln(os.Stdout, "Published message:", string(message))
		token.Wait()

		<-ticker.C
	}
}

func createJSONMessage(config Configuration, roundedValue float64) []byte {
	data := Data{
		Value:            roundedValue,
		Unit:             config.Unit,
		TransmissionRate: config.TransmissionRate,
		Longitude:        config.Longitude,
		Latitude:         config.Latitude,
		Sensor:           config.Sensor,
		Timestamp:        time.Now(),
		QoS:			  config.QoS,
	}

	jsonMsg, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return jsonMsg
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
