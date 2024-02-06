package main

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1891")
	opts.SetClientID("go_publisher")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		text := "Hello MQTT " + time.Now().Format(time.RFC3339)
		token := client.Publish("test/topic", 0, false, text)
		token.Wait()
		fmt.Println("Publicado:", text)
		time.Sleep(2 * time.Second)
	}
}