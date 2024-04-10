import pytest
import json
import time
import math
import paho.mqtt.client as paho
from paho import mqtt
from confluent_kafka import Producer, Consumer, KafkaError


mock_config = {
    "Sensor": "solar",
    "Longitude": 59.0,
    "Latitude": 55.0,
    "TransmissionRate": 10,
    "Unit": "W/m³",
    "QoS": 1,
}

mock_data = [2.25, 2.50, 4.25]
test_message = json.dumps({"Value": 1.25, "Unit": "W/m³", "TransmissionRate": 10, "Longitude": 59.0, "Latitude": 55.0, "Sensor": "air", "Timestamp": time.time(), "QoS": 1})
received_messages = []
first_message_timestamp = None
last_message_timestamp = None
received_qos = []
broker_address = 'b9f3c31144f64d469f184727678d8fb6.s1.eu.hivemq.cloud'
port = 8883
topic = "hi"
username = 'admin'
password = 'Admin123'
connected = False



def on_connect(client, userdata, flags, reason_code, properties):
    global connected
    if reason_code == 0:
        connected = True
    else:
        connected = False

def on_message(client, userdata, message):
    global received_messages, received_qos, first_message_timestamp, last_message_timestamp
    payload = str(message.payload.decode("utf-8"))
    print('here')
    received_messages.append(payload)
    received_qos.append(message.qos)

    if len(received_messages) == 1:
        first_message_timestamp = time.time()

    last_message_timestamp = time.time()

@pytest.fixture
def mqtt_client():
    global received_messages
    received_messages = []
    client = paho.Client(paho.CallbackAPIVersion.VERSION2, "test",
                     protocol=paho.MQTTv5)
    client.on_connect = on_connect

    # Configurações de TLS
    client.tls_set(tls_version=mqtt.client.ssl.PROTOCOL_TLS)
    client.username_pw_set(username, password)  # Configuração da autenticação

    # Conexão ao broker
    client.connect(broker_address, port=port)

    time_per_message = 1 / mock_config["TransmissionRate"]
    client.on_message = on_message

    client.subscribe(f"sensor/{mock_config['Sensor']}", qos=1)
    client.loop_start()

    return client


def test_mqtt_connection(mqtt_client):
    time.sleep(1)
    assert connected
    

def on_message(client, userdata, message):
    global received_messages, received_qos, first_message_timestamp, last_message_timestamp

    payload = json.loads(str(message.payload.decode("utf-8")))['Value']
    
    received_messages.append(payload)
    print("received messages", received_messages)
    received_qos.append(message.qos)

    if len(received_messages) == 1:
        first_message_timestamp = time.time()

    last_message_timestamp = time.time()


def test_message_reception(mqtt_client):

    global received_messages
    received_messages = []

    num_messages = len(mock_data)
    time_per_message = 1 / mock_config["TransmissionRate"]

    for data in mock_data:
        time.sleep(time_per_message)

        payload = json.dumps({"Value": data, "Unit": mock_config["Unit"], "TransmissionRate": mock_config["TransmissionRate"], "Longitude": mock_config["Longitude"], "Latitude": mock_config["Latitude"], "Sensor": mock_config["Sensor"], "Timestamp": time.time(), "QoS": mock_config["QoS"]})

        mqtt_client.publish(f"sensor/{mock_config['Sensor']}", payload=payload, qos=1)

    
    total_time = num_messages * time_per_message + 0.5
    time.sleep(total_time)

    assert len(received_messages) > 0

def test_kafka_integration(mqtt_client):
    

    consumer_config = {
    "bootstrap.servers":"pkc-rgm37.us-west-2.aws.confluent.cloud:9092",
    'security.protocol':"SASL_SSL",
    'sasl.mechanisms':'PLAIN',
    'sasl.username':'VSCWFUOLUOZZ2GME',
    'sasl.password':'jjdma8p/cWZks17T2fh9ZyDPKNsOf297oWvRm+RvIHg++J52/uYc6LfIrq/mF5Zq',
    'group.id': 'python-consumer-group',
    'auto.offset.reset': 'earliest'

}
    topic = 'sensor'
    consumer = Consumer(**consumer_config)

    consumer.subscribe([topic])

    time.sleep(2)

    print('publishing message')
    print(test_message)
    mqtt_client.publish(f"sensor/solar", payload=test_message, qos=1)

    # Consumir mensagens
    try:
        while True:
            msg = consumer.poll(timeout=1.0)
            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    continue
                else:
                    print(msg.error())
                    break
            print(f'Received message: {msg.value().decode("utf-8")}')
            received_value = json.loads(msg.value().decode("utf-8"))['Value']
            desired_value = json.loads(test_message)['Value']
            print(received_value, desired_value)
            if received_value == desired_value:
                assert True
                break
                

    except KeyboardInterrupt:
        pass
    finally:
        # Fechar consumidor
        consumer.close()

def test_mongodb_integration(mqtt_client):
    from pymongo import MongoClient
    from pymongo.server_api import ServerApi

    mqtt_client.publish(f"sensor/solar", payload=test_message, qos=1)
    uri = "mongodb+srv://admin:admin@sensor.fzgi4a5.mongodb.net/?retryWrites=true&w=majority&appName=sensor"
# Create a new client and connect to the server
    client = MongoClient(uri, server_api=ServerApi('1'))
    db = client['sensor']
    collection = db['reading']
    # find document by timestamp
    doc = collection.find().sort('_id',-1).limit(1)
    doc = list(doc)[0]
    print(doc)
    assert doc['Value'] == json.loads(test_message)['Value']