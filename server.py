from subprocess import call

from flask import Flask, request
from gevent.pywsgi import WSGIServer

app = Flask(__name__)

# TODO:
# Fix /ts for timestamp/hash
# Create js script to be added to the response to check for updates


@app.route('/', defaults={'path': ''}, methods=['GET'])
def get(path):
    try:
        with open('out', 'rb') as f:
            return f.read(), 200
    except Exception as e:
        return '', 500


@app.route('/', defaults={'path': ''}, methods=['POST'])
def post(path):
    try:
        with open('in', 'wb') as f:
            f.write(request.get_data())
        call(['pandoc', '-t', 'html', '-o', 'out', 'in'])
    except Exception as e:
        return '', 400
    return '', 200


@app.route('/ts', defaults={'path': ''}, methods=['GET'])
def ts(path):
    return '', 200


if __name__ == '__main__':
    http_server = WSGIServer(('', 8080), app)
    http_server.serve_forever()
