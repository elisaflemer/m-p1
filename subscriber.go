package main

import (
	"fmt"
	"bytes"
	"flag"
	"net/http"
	//"github.com/go-chi/chi"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)


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
	fmt.Printf("Recebido: %s do tópico: %s\n", msg.Payload(), msg.Topic())

	postStructAsJSON("http://localhost:5000/api", msg.Payload())

}

func subscribe() {

	configPath := flag.String("config", "", "Path to the configuration file")
	connection := flag.String("connection", "hivemq", "Enter 'hivemq' or 'local' for MQTT connection")

	hivemqUsername := flag.String("username", "", "HiveMQ username")
	hivemqPassword := flag.String("password", "", "HiveMQ password")

	flag.Parse()


	if *configPath == "" {
		fmt.Println("Usage: go run subscriber.go -config <config_path> -connection <hivemq/local> -username <username> -password <password>")
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


	client := connector.Connect("subscriber", *hivemqUsername, *hivemqPassword)
	defer client.Disconnect(250)

	if token := client.Subscribe("sensor/"+config.Sensor, 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	fmt.Println("Subscriber está rodando. Pressione CTRL+C para sair.")
	select {} // Bloqueia indefinidamente
}