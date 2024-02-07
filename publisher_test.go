package main

import (
	"fmt"
	"math"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var receivedMessages []string
var firstMessageTimestamp time.Time
var lastMessageTimestamp time.Time

var messagePubTestHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	receivedMessages = append(receivedMessages, payload)

	// Capture timestamps for the first and last messages
	if len(receivedMessages) == 1 {
		firstMessageTimestamp = time.Now()
	}

	lastMessageTimestamp = time.Now()
}

// Test for integrity
func TestIntegrity(t *testing.T) {
	client := connectMQTT("subscriber")
	defer client.Disconnect(250)

	mockConfig := Configuration{
		Sensor:           "air",
		Longitude:        59.0,
		Latitude:         55.0,
		TransmissionRate: 10,
		Unit:             "W/m³",
	}

	if token := client.Subscribe("sensor/"+mockConfig.Sensor, 1, messagePubTestHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	mockData := []float64{1.25, 2.50, 1.25, 2.50, 1.25, 2.50, 0, 0, 2.50, 1.25, 2.50}
	publishData(client, mockConfig, mockData)

	// Wait for a while to ensure the subscriber has received the message
	time.Sleep(5 * time.Second)

	// Perform integrity assertions or checks here
	if len(receivedMessages) == 0 {
		t.Errorf("No messages received")
	}

	if len(receivedMessages) != len(mockData) {
		t.Errorf("Received %d messages, expected %d", len(receivedMessages), len(mockData))
	}

	// Print timestamps
	fmt.Println("First Message Timestamp:", firstMessageTimestamp)
	fmt.Println("Last Message Timestamp:", lastMessageTimestamp)
}

// Test for transmission rate
func TestTransmissionRate(t *testing.T) {
	client := connectMQTT("subscriber")
	defer client.Disconnect(250)

	mockConfig := Configuration{
		Sensor:           "air",
		Longitude:        59.0,
		Latitude:         55.0,
		TransmissionRate: 10,
		Unit:             "W/m³",
	}

	if token := client.Subscribe("sensor/"+mockConfig.Sensor, 1, messagePubTestHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	mockData := []float64{1.25, 2.50, 1.25, 2.50, 1.25, 2.50, 0, 0, 2.50, 1.25, 2.50}
	publishData(client, mockConfig, mockData)

	// Wait for a while to ensure the subscriber has received the message
	time.Sleep(5 * time.Second)

	// Calculate time period in seconds
	timePeriod := lastMessageTimestamp.Sub(firstMessageTimestamp).Seconds()

	// Calculate frequency in Hz
	frequency := float64(len(mockData)) / timePeriod

	// Perform transmission rate assertions or checks here
	if math.Abs(frequency-mockConfig.TransmissionRate) > 2 {
		t.Errorf("Received frequency: %f, expected: %f", frequency, mockConfig.TransmissionRate)
	}
}
