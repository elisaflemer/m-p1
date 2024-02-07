# Ponderada 1, módulo 9

Nesta atividade, criamos um nó publisher em MQTT para valores simulados de um sensor. Garantimos a reutilização desse código através de abstrações como um arquivo em .json para as configurações de taxa de transmissão, tipo de sensor, latitude, longitude e unidade, e também um CSV com valores simulados de um sensor. Esses arquivos são passados como argumentos na linha de comando. Já para gerar os valores do CSV, utilizamos um script "generator.py", que recebe o número de dados a serem criados, o valor mínimo, o valor máximo e a resolução também pela liha de comando.

Os valores serão publicados como json com metadados no tópico "sensor/<nome-do-sensor>".

## Como rodar
Primeiro, é necessário gerar os dados para simulação. Para isso, execute os seguintes comandos neste diretório:

```
pip install csv
```

```
python3 generator.py <num_de_valores> <resolucao> <valor_minimo> <valor_maximo>
```

Isso gerará um CSV denominado "data.csv".

A partir daí, podemos simular a leitura desse CSV (ou de qualquer outro). Para isso, será necessário criar também um arquivo de configuração em json, no seguinte padrão:

```
{
    "sensor": <nome_do_sensor>,
    "longitude": <longitude>,
    "latitude": <latitude>,
    "transmission_rate_hz": <transmission_rate_in_herz>,
    "unit": <unit>
}
```

Feito isso, execute:
```
go mod tidy
go run publisher.go <config_path> <csv_path>
```

## Demo
[demo.webm](https://github.com/elisaflemer/m-p1/assets/99259251/4328ea4a-ce58-467c-b49a-04c3e6ec1c34)
