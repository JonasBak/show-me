from subprocess import call

import hashlib

from flask import Flask, request
from gevent.pywsgi import WSGIServer

app = Flask(__name__)

fh = ''
with open('reload.html', 'r') as f:
    reload = f.read()


@app.route('/', defaults={'path': ''}, methods=['GET'])
def get(path):
    try:
        with open('out', 'r') as f:
            file = f.read()
        return reload.replace('/*hash*/', f'"{fh}"') + '\n' + file, 200
    except Exception as e:
        return 'No file loaded yet', 404


@app.route('/', defaults={'path': ''}, methods=['POST'])
def post(path):
    global fh
    with open('in', 'wb') as f:
        file = request.get_data()
        fh = hashlib.md5(file).hexdigest()
        f.write(file)
    call(['pandoc', '-t', 'html', '-o', 'out', 'in'])
    return '', 200


@app.route('/ts', defaults={'path': ''}, methods=['GET'])
def ts(path):
    return fh, 200


if __name__ == '__main__':
    http_server = WSGIServer(('', 8080), app)
    http_server.serve_forever()
