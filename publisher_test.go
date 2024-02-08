package main

import (
	"math"
	"testing"
	"time"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var receivedMessages []string
var firstMessageTimestamp time.Time
var lastMessageTimestamp time.Time
var receivedQoS []byte

var messagePubTestHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	payload := string(msg.Payload())
	receivedMessages = append(receivedMessages, payload)

	// Capture QoS
	receivedQoS = append(receivedQoS, msg.Qos())

	// Capture timestamps for the first and last messages
	if len(receivedMessages) == 1 {
		firstMessageTimestamp = time.Now()
	}

	lastMessageTimestamp = time.Now()
}

func TestIntegration(t *testing.T) {
	client := connectMQTT("subscriber")
	defer client.Disconnect(250)

	mockConfig := Configuration{
		Sensor:           "air",
		Longitude:        59.0,
		Latitude:         55.0,
		TransmissionRate: 10,
		Unit:             "W/m³",
		QoS:			  1,
	}

	if token := client.Subscribe("sensor/"+mockConfig.Sensor, 1, messagePubTestHandler); token.Wait() && token.Error() != nil {
		t.Fatalf("Error subscribing to MQTT: %s", token.Error())
	}

	mockData := []float64{1.25, 2.50, 1.25, 2.50, 1.25, 2.50, 0, 0, 2.50, 1.25, 2.50}
	publishData(client, mockConfig, mockData)

	// Wait for a while to ensure the subscriber has received the message
	time.Sleep(5 * time.Second)

	// Check integrity
	if len(receivedMessages) == 0 {
		t.Fatal("\x1b[31m[FAIL] No messages received\x1b[0m")
	} else {
		t.Log("\x1b[32m[PASS] Messages received successfully\x1b[0m")
	}

	if len(receivedMessages) != len(mockData) {
		t.Fatalf("\x1b[31m[FAIL] Received %d messages, expected %d\x1b[0m", len(receivedMessages), len(mockData))
	} else {
		t.Log("\x1b[32m[PASS] Correct number of messages received\x1b[0m")
	}

	// Check QoS
	QoSFail := false
	for i, qos := range receivedQoS {
		if qos != mockConfig.QoS {
			t.Fatalf("\x1b[31m[FAIL] Incorrect QoS in message %d. Received QoS: %d, expected: 1\x1b[0m", i, qos)
			QoSFail = true
			} 
		}

	if (!QoSFail) {
		t.Log("\x1b[32m[PASS] Correct QoS received\x1b[0m")
	}

	// Decode JSON messages
	var decodedMessages []float64
	for _, msg := range receivedMessages {
		var m Data
		if err := json.Unmarshal([]byte(msg), &m); err != nil {
			t.Fatalf("Error decoding JSON: %s", err)
		}
		decodedMessages = append(decodedMessages, m.Value)
	}

	// Check if received messages match expected values
	if fmt.Sprintf("%v", decodedMessages) != fmt.Sprintf("%v", mockData) {
		t.Fatalf("\x1b[31m[FAIL] Received %v, expected %v\x1b[0m", decodedMessages, mockData)
	} else {
		t.Log("\x1b[32m[PASS] Correct messages received\x1b[0m")
	}

	// Calculate time period in seconds
	timePeriod := lastMessageTimestamp.Sub(firstMessageTimestamp).Seconds()

	// Calculate frequency in Hz
	frequency := float64(len(mockData)) / timePeriod

	// Check transmission rate
	if math.Abs(frequency-mockConfig.TransmissionRate) > 2 {
		t.Fatalf("\x1b[31m[FAIL] Received frequency: %f, expected: %f\x1b[0m", frequency, mockConfig.TransmissionRate)
	} else {
		t.Log("\x1b[32m[PASS] Transmission rate within acceptable range of 2Hz\x1b[0m")
	}
}
