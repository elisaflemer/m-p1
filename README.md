# Ponderada de integração

Este documento fornece uma visão geral do fluxo de dados e das conexões entre um publicador de dados em Python e um banco de dados MongoDB, intermediado por um broker MQTT (HiveMQ) e um sistema de mensagens (Kafka). O código disponibilizado faz a ponte entre o dispositivo sensorial (publicador) e o banco de dados MongoDB, permitindo a captura, processamento e armazenamento dos dados provenientes do sensor.

## Visão Geral do Fluxo de Dados
**Publicador de Dados em Python:** O código Python disponibilizado representa o publicador de dados. Ele simula um sensor que coleta informações, como valor, unidade, taxa de transmissão, coordenadas geográficas, sensor utilizado e carimbo de data/hora. Esses dados são formatados em JSON e publicados em um tópico MQTT específico.

**Broker MQTT (HiveMQ):** O HiveMQ é o broker MQTT utilizado para facilitar a comunicação entre o publicador Python e o sistema de mensagens (Kafka). Ele recebe os dados do publicador e os encaminha para o Kafka, utilizando um tópico MQTT.

**Sistema de Mensagens (Kafka):** O Kafka é um sistema de mensagens distribuído que recebe os dados do broker MQTT (HiveMQ) e os armazena em tópicos específicos. Ele atua como um intermediário entre o HiveMQ e o MongoDB.

**Banco de Dados MongoDB:** O MongoDB é um banco de dados NoSQL utilizado para armazenar os dados recebidos do Kafka. Os dados são processados e persistidos no MongoDB para futuras consultas e análises.

## Configuração e Execução

### Requisitos de Software

- Python 3.x
- Bibliotecas Python: paho-mqtt, confluent-kafka, pytest, pymongo

### Configuração do Ambiente

Execute o seguinte comando para instalar as bibliotecas Python necessárias:

```
pip install paho-mqtt confluent-kafka pytest pymongo
```

Certifique-se de ter um servidor MongoDB em execução e substitua as credenciais e detalhes de conexão no código de teste (test_mongodb_integration) conforme necessário.

Execute o script publisher.py para iniciar o publicador de dados. Certifique-se de ter o arquivo de configuração config.json e o arquivo CSV de dados data.csv no diretório de trabalho.

```
python publisher.py
```

## Execução dos Testes
Execute os testes unitários usando o pytest. Certifique-se de ter os servidores HiveMQ, Kafka e MongoDB em execução antes de executar os testes.

```
pytest -v
```

### Detalhes do Código
- publisher.py: Contém o código do publicador de dados em Python.
- test_publisher.py: Contém os testes unitários para validar a conexão com o broker MQTT, recepção de mensagens MQTT e integração com Kafka e MongoDB.
- config.json: Arquivo de configuração com os detalhes do sensor e configurações de conexão.
- data.csv: Arquivo CSV contendo os dados simulados do sensor.

### Testes Automatizados
Os testes automatizados validam diferentes aspectos do fluxo de dados, incluindo a conexão MQTT, recepção de mensagens MQTT, integração com Kafka e integração com MongoDB. Eles garantem que o sistema esteja funcionando corretamente e que os dados sejam transmitidos e armazenados conforme o esperado.