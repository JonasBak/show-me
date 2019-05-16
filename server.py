import hashlib
import sys
from subprocess import call

from flask import Flask, request
from gevent.pywsgi import WSGIServer

app = Flask(__name__)

fh = ''


@app.route('/', defaults={'path': ''}, methods=['GET'])
def get(path):
    try:
        with open('out', 'r') as f:
            file = f.read()
        return file.replace('/*hash*/', f'"{fh}"'), 200
    except Exception as e:
        return 'No file loaded yet', 404


@app.route('/', defaults={'path': ''}, methods=['POST'])
def post(path):
    global fh
    with open('in', 'wb') as f:
        file = request.get_data()
        fh = hashlib.md5(file).hexdigest()
        f.write(file)
    args = [
        'pandoc', '-t', 'html', '-o', 'out', '-H', 'reload.html', '-M',
        f'title={request.args.get("filename", " ")}'
    ]
    if (len(request.args.get('from', '')) > 0):
        args += ['-f', request.args.get('from')]
    args += ['in']
    call(args)
    return '', 200


@app.route('/ts', defaults={'path': ''}, methods=['GET'])
def ts(path):
    return fh, 200


@app.after_request
def flush_streams(response):
    sys.stdout.flush()
    sys.stderr.flush()
    return response


if __name__ == '__main__':
    http_server = WSGIServer(('', 8080), app)
    http_server.serve_forever()
