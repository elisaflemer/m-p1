import csv
import random
# Função para gerar leituras de radiação solar
def gerar_leitura():
    # Resolução de 1,25
    valor = round(random.uniform(0, 1280), 2)
    return valor

# Número de leituras a serem geradas
num_leituras = 100

# Gerar leituras e escrever no arquivo CSV
with open("leituras_solar.csv", "w", newline="") as csvfile:
    writer = csv.writer(csvfile)
    
    # Gerar e escrever leituras
    for _ in range(num_leituras):
        valor = gerar_leitura()
        writer.writerow([valor])

print(f"Arquivo 'leituras_solar.csv' gerado com sucesso.")
