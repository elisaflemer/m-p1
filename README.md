# Ponderada 1, módulo 9

Nesta atividade, criamos um nó publisher em MQTT para valores simulados de um sensor. Garantimos a reutilização desse código através de abstrações como um arquivo em .json para as configurações de taxa de transmissão, tipo de sensor, latitude, longitude e unidade, e também um CSV com valores simulados de um sensor. Esses arquivos são passados como argumentos na linha de comando. Já para gerar os valores do CSV, utilizamos um script "generator.py", que recebe o número de dados a serem criados, o valor mínimo, o valor máximo e a resolução também pela liha de comando.

Os valores são publicados como json com metadados no tópico "sensor/<nome-do-sensor>".

Todas as principais funções, como conexão com o broken, integridade de mensagens, taxa de transmissão e QoS são testadas com testes automáticos em Go.

## Atualização para ponderada 4

Para a ponderada 4, o código foi atualizado para que se conecte com o broker da HiveMQ mediante autenticação definida na plataforma do cluster. O teste de transmissão também foi atualizado para enviar 1000 mensagens com taxa de 10 Hz, tendo passado dentro e uma margem de erro de 2Hz.

## Como rodar

### 1. Gerar dados de simulação

Primeiro, é necessário gerar os dados para simulação. Para isso, execute os seguintes comandos neste diretório:

```
pip install csv
```

```
python3 generator.py <num_de_valores> <resolucao> <valor_minimo> <valor_maximo>
```

Isso gerará um CSV denominado "data.csv".

A partir daí, podemos simular a leitura desse CSV (ou de qualquer outro). 

### 2. Criar arquivo de configuração

Crie um arquivo de configuração JSON com informações sobre o sensor. O arquivo deve seguir o seguinte padrão:
```
{
    "sensor": <nome_do_sensor>,
    "longitude": <longitude>,
    "latitude": <latitude>,
    "transmission_rate_hz": <taxa_de_transmissao_em_herz>,
    "unit": <unidade>,
    "QoS": <qos>
}
```

### 3. Executar o Publisher MQTT

Certifique-se de ter o mosquitto instalado com arquivo de configuração (ouvindo na porta 1891) e inicie o broker local:

```
mosquitto -c mosquitto.conf
```

No diretório deste projeto, instale dependências do Go:

```
go mod tidy
```

Execute o script publisher.go passando o caminho do arquivo de configuração JSON e o caminho do arquivo CSV:

```
go run publisher.go <config_path> <csv_path>
```

## Estrutura dos dados

### Configuração (arquivo json)
- unit: Unidade do sensor.
- transmission_rate_hz: Taxa de transmissão em Hertz.
- longitude: Longitude do sensor.
- latitude: Latitude do sensor.
- sensor: Nome do sensor.
- qos: Qualidade de serviço para comunicação MQTT.

### Dados do Sensor (json publicado no tópico)
- value: Valor lido pelo sensor (proveniente do CSV)
- unit: Unidade do sensor.
- transmission_rate_hz: Taxa de transmissão em Hertz.
- longitude: Longitude do sensor.
- latitude: Latitude do sensor.
- sensor: Nome do sensor.
- timestamp: Timestamp da leitura do sensor.
- qos: Qualidade de serviço para comunicação MQTT.

## Como testar

O projeto inclui testes unitários para garantir a robustez e a integridade do código. Os testes são implementados usando o pacote nativo testing do Go e abrangem as seguintes áreas-chave:

### 1. Conexão com broker

Este teste verifica se é possível conectar-se com sucesso ao broker MQTT. Ele utiliza a função connectMQTT para estabelecer uma conexão e, em seguida, verifica se a conexão foi bem-sucedida usando client.IsConnected(). Se a conexão for estabelecida com sucesso, o teste é considerado passado; caso contrário, é considerado falho.

```go
func TestConnectMQTT(t *testing.T) {
    // Conectar ao broker MQTT
    client := connectMQTT("publisher")
    defer client.Disconnect(250)

    // Verificar se a conexão foi estabelecida com sucesso
    if !client.IsConnected() {
        t.Fatalf("\x1b[31m[FAIL] Unable to connect to MQTT broker\x1b[0m")
    } else {
        t.Log("\x1b[32m[PASS] Connected to MQTT broker\x1b[0m")
    }
}
```
### 2. Chegada de mensagens

Este teste verifica se as mensagens são recebidas corretamente. Ele configura um ambiente de teste usando a função setupTest, que cria uma subscrição MQTT para o tópico específico e publica dados simulados. Em seguida, o teste aguarda um período de tempo que permitiria a recepção de todas as mensagens esperadas. Se nenhuma mensagem for recebida, o teste falha; caso contrário, é considerado passado. Utilizamos essa condição porque, dependendo do QoS escolhido, pode haver o recebimento de mais mensagens do que enviadas; portanto, o a comparação entre quantas mensagens foram enviadas e recebidas não é um bom critério. Em vez disso, checamos a integridade das mensagens no teste seguinte.

```go
func TestMessageReception(t *testing.T) {
    // Configurar o ambiente de teste
    setupTest(t)

    // Aguardar o tempo necessário para a recepção das mensagens
    timePerMessage := time.Duration(int(time.Second) / int(mockConfig.TransmissionRate))
    timeMargin := int(0.5 * float64(time.Second))
    totalTime := time.Duration(len(mockData)*int(timePerMessage)+timeMargin)
    time.Sleep(totalTime)

    // Verificar se pelo menos uma mensagem foi recebida
    if len(receivedMessages) == 0 {
        t.Fatal("\x1b[31m[FAIL] No messages received\x1b[0m")
    } else {
        t.Log("\x1b[32m[PASS] Messages received successfully\x1b[0m")
    }
}
```

### 3. Integridade das Mensagens

O teste `TestMessageIntegrity` garante a integridade das mensagens recebidas. Ele decodifica cada mensagem recebida e verifica se cada valor em `mockData` tem pelo menos uma correspondência em `decodedMessages`. Isso garante que os valores esperados estejam presentes nas mensagens recebidas, embora não necessariamente na mesma ordem.

```go
func TestMessageIntegrity(t *testing.T) {
    // Configurar o ambiente de teste
    setupTest(t)

    // Decodificar mensagens e verificar integridade
    var decodedMessages []float64
    for _, msg := range receivedMessages {
        var m Data
        if err := json.Unmarshal([]byte(msg), &m); err != nil {
            t.Fatalf("Error decoding JSON: %s", err)
        }
        decodedMessages = append(decodedMessages, m.Value)
    }

    // Verificar se cada valor em mockData tem pelo menos uma correspondência em decodedMessages
    for _, expectedValue := range mockData {
        found := false
        for _, decodedValue := range decodedMessages {
            if expectedValue == decodedValue {
                found = true
                break
            }
        }
        if !found {
            t.Fatalf("\x1b[31m[FAIL] Value %v not found in received messages: %v\x1b[0m", expectedValue, decodedMessages)
        }
    }

    t.Log("\x1b[32m[PASS] Correct messages received\x1b[0m")
}
```

Este teste assegura que cada valor em mockData seja incluído nas mensagens recebidas, garantindo assim a integridade e a correspondência dos dados.

### 4. Taxa de transmissão

Este teste verifica se a taxa de transmissão das mensagens está dentro de uma faixa aceitável (default de +/- 2Hz) em relação à taxa configurada. Ele calcula o período de tempo entre a primeira e a última mensagem recebida e, com base nisso, calcula a frequência real das mensagens. Se a frequência estiver fora da faixa aceitável, o teste falha; caso contrário, é considerado passado.

```go
func TestTransmissionRate(t *testing.T) {
    // Configurar o ambiente de teste
    setupTest(t)

    // Calcular o período de tempo e a frequência
    timePeriod := lastMessageTimestamp.Sub(firstMessageTimestamp).Seconds()
    frequency := float64(len(mockData)) / timePeriod

    // Verificar se a taxa de transmissão está dentro da faixa aceitável
    if math.Abs(frequency-mockConfig.TransmissionRate) > 2 {
        t.Fatalf("\x1b[31m[FAIL] Received frequency: %f, expected: %f\x1b[0m", frequency, mockConfig.TransmissionRate)
    } else {
        t.Log("\x1b[32m[PASS] Transmission rate within acceptable range of 2Hz\x1b[0m")
    }
}

```

### 5. QoS

Este teste avalia a correta entrega de mensagens conforme configurado pelo QoS. Ele publica uma única mensagem simulada com a configuração de QoS especificada, verifica se ela foi entregue corretamente segundo seu QoS e, em seguida, relata o resultado do teste.

```go
func TestQoS(t *testing.T) {
	client := connectMQTT("subscriber")
	defer client.Disconnect(250)

	if token := client.Subscribe("sensor/"+mockConfig.Sensor, mockConfig.QoS, messagePubTestHandler); token.Wait() && token.Error() != nil {
		t.Fatalf("Error subscribing to MQTT: %s", token.Error())
	}
	receivedMessages = []string{}
	mockQoSData := []float64{1.25}
	publishData(client, mockConfig, mockQoSData)
	time.Sleep(time.Duration(1 / int(mockConfig.TransmissionRate) * int(time.Second)))

	switch mockConfig.QoS {
	case 0:
		t.Log("\x1b[33m[INFO] QoS set to 0, no guarantee of message delivery\x1b[0m")
	case 1:
		if len(receivedMessages) == 0 {
			t.Fatalf("\x1b[31m[FAIL] No messages received with QoS 1\x1b[0m")
		} else {
			for _, msg := range receivedMessages {
				var m Data
				if err := json.Unmarshal([]byte(msg), &m); err != nil {
					t.Fatalf("Error decoding JSON: %s", err)
				}
				if m.Value != mockQoSData[0] {
					t.Fatalf("\x1b[31m[FAIL] Received %v, expected %v\x1b[0m", m.Value, mockQoSData[0])
				}
			}
			t.Log("\x1b[32m[PASS] Message received with QoS 1\x1b[0m")
		}
	case 2:
		if len(receivedMessages) != 1 {
			t.Fatalf("\x1b[31m[FAIL] Incorrect number of messages received with QoS 2. Expected: 1, received: %d\x1b[0m", len(receivedMessages))
		} else {
			var m Data
				if err := json.Unmarshal([]byte(receivedMessages[0]), &m); err != nil {
					t.Fatalf("Error decoding JSON: %s", err)
				}
				if m.Value != mockQoSData[0] {
					t.Fatalf("\x1b[31m[FAIL] Received %v, expected %v\x1b[0m", m.Value, mockQoSData[0])
				}
				t.Log("\x1b[32m[PASS] Message received with QoS 2\x1b[0m")
		}
		default:
		t.Fatalf("\x1b[31m[FAIL] Invalid QoS value: %d\x1b[0m", mockConfig.QoS)
	}

}

```

## Demo
[demo_completa.webm](https://github.com/elisaflemer/m9-p1/assets/99259251/cc8a6a14-5036-48f3-a703-de2b0408b011)

![image](https://github.com/elisaflemer/m9-p1/assets/99259251/2da0ef29-135c-45f7-85fa-5e8a885eee08)
