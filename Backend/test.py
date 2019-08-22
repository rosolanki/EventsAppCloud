from flask import Flask, request, jsonify, Response
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

        print(data)
        return jsonify(data)
    if request.method == 'GET':
        user = request.args.get('user')
        filename = 'data.json'
        with open(filename, 'r+') as json_file:
            data = json.load(json_file)
        for item in data['feeds']:
            if item['user'] == user:
                iteminfo = json.dumps(item)
                resp = Response(iteminfo, status=200, mimetype='application/json')
                return jsonify(item)
        errorJson = {}
        errorJson["Error"] = "UserNotFound"
        return jsonify(errorJson)
        
app.run(host="0.0.0.0", port=8999)


{"jOHN":"tHIS IS jOHN"}

mydata["john"]

this is john