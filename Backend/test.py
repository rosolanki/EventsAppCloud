from flask import Flask, request, jsonify
import json
import os 

app = Flask(__name__)
@app.route('/Events', methods = ['POST', 'GET'])
def eventsProcess():
    if request.method == 'POST':
        mydata = request.get_json()
        mydata['Manipulated'] = "Yes"
        filename = 'data.json'
        with open(filename, 'r+') as json_file:
            data = json.load(json_file)
            data['feeds'].append(mydata)
        os.remove(filename)
        with open(filename, 'w') as json_file:
            json.dump(data, json_file)
        return jsonify(data)
    if request.method == 'GET':
        print("Nothing")
        
app.run(host="0.0.0.0", port=8999)