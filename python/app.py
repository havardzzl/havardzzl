from flask import Flask, render_template

app = Flask(__name__)

@app.route('/hi/<name>')
def hi(name):
    return f'Hi, {name}!'

@app.route("/")
def index():
    return render_template('index.html')

if __name__ == '__main__':
    app.run(port=9000, host='0.0.0.0')