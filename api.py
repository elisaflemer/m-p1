from flask import Flask, request
import sqlite3
import json
from datetime import datetime

app = Flask(__name__)

# Initialize the database
def init_db():
    conn = sqlite3.connect('ponderada.db')
    c = conn.cursor()

    # Create the table if it doesn't exist
    c.execute('''CREATE TABLE IF NOT EXISTS sensor_data (
    timestamp TIMESTAMP PRIMARY KEY,
    value FLOAT NOT NULL,
    unit VARCHAR(255) NOT NULL,
    transmission_rate FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    latitude FLOAT NOT NULL,
    sensor VARCHAR(255) NOT NULL,
    qos VARCHAR(255) NOT NULL
);''')

    conn.commit()
    conn.close()

# API endpoint for receiving JSON data
@app.route('/api', methods=['POST'])
def receive_data():
    data = request.get_json()

    print(data)

    # Store the data in the database
    conn = sqlite3.connect('ponderada.db')
    c = conn.cursor()

    c.execute("INSERT INTO sensor_data VALUES (?, ?, ?, ?, ?, ?, ?, ?)", (data['timestamp'], data['value'], data['unit'], data['transmission_rate'], data['longitude'], data['latitude'], data['sensor'], data['qos']))

    conn.commit()
    conn.close()

    return 'Data stored successfully'

# API endpoint for querying the database
@app.route('/api', methods=['GET'])
def query_data():
    conn = sqlite3.connect('ponderada.db')
    c = conn.cursor()

    # Get the data from the database
    c.execute("SELECT * FROM sensor_data")
    data = c.fetchall()

    conn.close()

    # Convert the data to JSON
    data_json = []
    for row in data:
        data_json.append({
            'timestamp': row[0],
            'value': row[1],
            'unit': row[2],
            'transmission_rate': row[3],
            'longitude': row[4],
            'latitude': row[5],
            'sensor': row[6],
            'qos': row[7]
        })

    return json.dumps(data_json)

if __name__ == '__main__':
    init_db()
    app.run()