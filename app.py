import http.server
import socketserver
import re

PORT = 9000

class MyHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        # Use a regular expression to match against the request path and extract parameters
        match = re.match(r'^/hello/(\w+)$', self.path)
        if match:
            name = match.group(1)
            self.send_response(200)
            self.end_headers()
            self.wfile.write(f'Hello, {name}!'.encode())
        else:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'Hello')

with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()