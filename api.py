from flask import Flask, request
from pymongo import MongoClient
import json
from datetime import datetime

app = Flask(__name__)

# MongoDB connection string. Replace <username>, <password>, and <database_name> with your actual MongoDB credentials.
mongo_uri = "mongodb+srv://elisaflemer:elisaflemer@sensor.yd8n6bm.mongodb.net/?retryWrites=true&w=majority&appName=Sensor"
client = MongoClient(mongo_uri)

# Get the database
db = client['Sensor']

# Initialize the database (for MongoDB, you might not need this step as collections are created dynamically)
def init_db():
    # Ensure the collection exists by attempting to count documents in it.
    # MongoDB creates collections when the first document is inserted.
    if 'Sensor' in db.list_collection_names():
        print("Collection exists.")
    else:
        print("Collection does not exist, it will be created on first insert.")

# API endpoint for receiving JSON data
@app.route('/api', methods=['POST'])
def receive_data():
    print('hi')
    data = request.get_json()

    print(data)

    # Store the data in the database
    db.Sensor.insert_one(data)

    return 'Data stored successfully'

# API endpoint for querying the database
@app.route('/api', methods=['GET'])
def query_data():
    # Get the data from the database
    cursor = db.Sensor.find()

    # Convert the data to JSON
    data_json = []
    for document in cursor:
        # MongoDB stores timestamps in a different format, you may need to convert them
        document['_id'] = str(document['_id'])  # Convert ObjectId to string to make it JSON-serializable
        data_json.append(document)

    return json.dumps(data_json, default=str)  # `default=str` to handle datetime objects

if __name__ == '__main__':
    init_db()
    app.run()
