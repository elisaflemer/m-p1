package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"math"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Configuration struct {
	Unit             string  `json:"unit"`
	TransmissionRate float64 `json:"transmission_rate_hz"` // in Hz
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	Sensor           string  `json:"sensor"`
}

type Data struct {
	Value            float64   `json:"value"`
	Unit             string    `json:"unit"`
	TransmissionRate float64   `json:"transmission_rate"`
	Longitude        float64   `json:"longitude"`
	Latitude         float64   `json:"latitude"`
	Sensor           string    `json:"sensor"`
	Timestamp        time.Time `json:"timestamp"`
}

func main() {
	// Read configuration from config.json
	config, err := readConfig("config.json")
	if err != nil {
		panic(err)
	}

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_publisher")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	file, err := os.Open("leituras_solar.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	

	// Calculate the time interval based on the transmission rate
	interval := time.Second / time.Duration(config.TransmissionRate)

	// Create a ticker to trigger data publication at regular intervals
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Process the line and extract the value
		// For simplicity, let's assume the value is the first field in each line
		var value float64
		_, err := fmt.Sscanf(line, "%f", &value)
		if err != nil {
			panic(err)
		}

		roundedValue := math.Round(value*100)/100

		// Create a JSON message with timestamp
		data := Data{
			Value:            roundedValue,
			Unit:             config.Unit,
			TransmissionRate: config.TransmissionRate,
			Longitude:        config.Longitude,
			Latitude:         config.Latitude,
			Sensor:           config.Sensor,
			Timestamp:        time.Now(),
		}
		jsonMsg, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		// Publish the JSON message
		token := client.Publish("test/topic", 0, false, jsonMsg)
		token.Wait()
		fmt.Println("Published:", string(jsonMsg))

		// Wait for the next tick
		<-ticker.C
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Read configuration from config.json
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
