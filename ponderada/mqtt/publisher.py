import json
import csv
import math
import argparse
from datetime import datetime, timedelta
import paho.mqtt.client as mqtt
import time
import paho.mqtt.client as paho
from paho import mqtt


class Configuration:
    def __init__(self, unit, transmission_rate_hz, longitude, latitude, sensor, qos):
        self.unit = unit
        self.transmission_rate_hz = transmission_rate_hz
        self.longitude = longitude
        self.latitude = latitude
        self.sensor = sensor
        self.qos = qos

class Data:
    def __init__(self, value, unit, transmission_rate, longitude, latitude, sensor, timestamp, qos):
        self.value = value
        self.unit = unit
        self.transmission_rate = transmission_rate
        self.longitude = longitude
        self.latitude = latitude
        self.sensor = sensor
        self.timestamp = timestamp
        self.qos = qos


    
def on_connect(client, userdata, flags, reason_code, properties):
    print(f"CONNACK received with code {reason_code}")


def read_config(filename):
    with open(filename, 'r') as f:
        data = json.load(f)
        return Configuration(**data)

def read_csv(csv_path):
    values = []
    with open(csv_path, newline='') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            values.append(float(row[0]))
    return values

def publish_data(client, config, data):
    interval = timedelta(seconds=1/config.transmission_rate_hz)
    for value in data:
        rounded_value = round(value, 2)
        print(f"Publishing: {rounded_value}")
        message = create_json_message(config, rounded_value)
        client.publish(f"sensor/{config.sensor}", message, qos=config.qos)
        time.sleep(interval.total_seconds())

def create_json_message(config, rounded_value):
    data = Data(
        value=rounded_value,
        unit=config.unit,
        transmission_rate=config.transmission_rate_hz,
        longitude=config.longitude,
        latitude=config.latitude,
        sensor=config.sensor,
        timestamp=datetime.now().isoformat(),
        qos=config.qos
    )
    return json.dumps(data.__dict__)

def main():
    broker_address = 'b9f3c31144f64d469f184727678d8fb6.s1.eu.hivemq.cloud'
    port = 8883
    config = read_config('config.json')
    values = read_csv('data.csv')
    topic = f"sensor/{config.sensor}"
    username = 'elisa'
    password = 'Elisa123'

    client = paho.Client(paho.CallbackAPIVersion.VERSION2, "Publisher",
                     protocol=paho.MQTTv5)

    client.on_connect = on_connect

    # Configurações de TLS
    client.tls_set(tls_version=mqtt.client.ssl.PROTOCOL_TLS)
    client.username_pw_set(username, password)  # Configuração da autenticação

    # Conexão ao broker
    client.connect(broker_address, port=port)

    publish_data(client, config, values)

if __name__ == "__main__":
    main()
