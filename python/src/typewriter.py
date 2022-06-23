import flask
import markdown
import os
import yaml

app = flask.Flask(__name__)
page_data = {}

@app.route('/api/foobar')
def api():
    print(flask.request.headers)
    return '200'

@app.route('/', defaults={'path': ''})
@app.route('/<path:path>')
def render(path):
    print(flask.request.headers)
    if path in page_data.keys():
        resp = flask.Response(markdown.markdown(page_data[path]['contents'], extensions=['mdx_math']))
        resp.headers['X-Forwarded-For'] = 'foobar'
        return resp
    else:
        return '404'

def load_paths():
    for fi in os.listdir('markdown'):
        with open('markdown/' + fi) as f:
            header = ''
            contents = ''
            lines = f.read().split('\n')
            is_header = True
            for l in lines:
                if l.startswith('---'):
                    is_header = False
                    continue
                if is_header:
                    header += l + '\n'
                else:
                    contents += l + '\n'
            metadata = yaml.safe_load(header)
            key = metadata['path']
            page_data[key] = {'metadata': metadata, 'contents': contents}

if __name__ == '__main__':
    load_paths()
    print(page_data)
    app.run()
