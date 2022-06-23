import flask
import requests
from urllib.parse import urlparse

app = flask.Flask(__name__)
proxies = {'proxy':'http://localhost:5000'}

@app.route('/', defaults={'path': ''})
@app.route('/<path:path>')
def proxy(path):
    print(flask.request.headers)
    for proxy_path, server in proxies.items():
        if path.startswith(proxy_path):
            proxy_count = len(proxy_path.split('/'))
            new_path = '/'.join(path.split('/')[proxy_count:])
            response = requests.get(f'{server}/{new_path}')
            return response.content
        referrer = flask.request.referrer
        if not referrer is None:
            referrer_path = urlparse(referrer).path[1:]
            if referrer_path.startswith(proxy_path):
                return flask.redirect(f"/{proxy_path}/{path}", code=302)
    return '404', 404



if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)
