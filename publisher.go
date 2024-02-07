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
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run publisher.go <config_path> <csv_path>")
		return
	}

	// Extract command-line arguments
	configPath := os.Args[1]
	csvPath := os.Args[2]

	config, err := readConfig(configPath)
	if err != nil {
		panic(err)
	}

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_publisher")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	file, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	interval := time.Second / time.Duration(config.TransmissionRate)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

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
		token := client.Publish("sensor/" + config.Sensor, 0, false, jsonMsg)
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
