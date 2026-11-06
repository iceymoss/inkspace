"""Exercise the operator smoke script against a disposable HTTP fixture."""

import http.server
import json
import subprocess
import threading
import unittest
from pathlib import Path


BASE = "https://canonical.example"
ARTICLE = BASE + "/blog/public-fixture"
PRIVATE = "hidden-fixture"


class Fixture(http.server.BaseHTTPRequestHandler):
    bad_canonical = False

    def log_message(self, *_args):
        pass  # never log fixture paths or hostile request headers

    def do_HEAD(self):
        self.do_GET()

    def do_GET(self):
        path = self.path.split("?")[0]
        headers = {}
        status = 200
        mime = "text/html; charset=utf-8"
        canonical = "https://incorrect.example/blog/public-fixture" if self.bad_canonical else ARTICLE
        body = ""
        if path == "/blog/1":
            status, headers["Location"] = 301, ARTICLE
        elif path == "/blog/" + PRIVATE:
            status = 404
        elif path == "/blog/public-fixture":
            headers["Cache-Control"] = "no-cache"
            body = ('<link rel="canonical" href="' + canonical + '">'
                    '<meta property="og:url" content="' + ARTICLE + '">'
                    '<meta name="description" content="公开摘要">'
                    '<meta name="twitter:card" content="summary">'
                    '<link rel="alternate" href="' + BASE + '/feed.xml">'
                    '<link rel="alternate" href="' + BASE + '/atom.xml">'
                    '<script type="application/ld+json">' +
                    json.dumps({"@type": "BlogPosting", "url": ARTICLE}) + '</script>')
        elif path == "/api/settings/public":
            body = json.dumps({"code": 0, "data": {"public_site_url": BASE}})
        else:
            headers.update({"ETag": '"fixture"', "Cache-Control": "public, max-age=300"})
            if path == "/feed.xml":
                mime = "application/rss+xml; charset=utf-8"
                body = '<rss><channel><link>' + BASE + '/</link><item><guid>' + ARTICLE + '</guid></item></channel></rss>'
            elif path == "/atom.xml":
                mime = "application/atom+xml; charset=utf-8"
                body = '<feed xmlns="http://www.w3.org/2005/Atom"><id>' + BASE + '/</id><entry><id>' + ARTICLE + '</id></entry></feed>'
            elif path in ("/sitemap.xml", "/sitemaps/1.xml"):
                mime = "application/xml; charset=utf-8"
                body = '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>' + ARTICLE + '</loc></url></urlset>'
            elif path == "/robots.txt":
                mime, body = "text/plain; charset=utf-8", "Sitemap: " + BASE + "/sitemap.xml"
            else:
                status = 404
            if self.headers.get("If-None-Match") == 'W/"fixture"':
                status, body = 304, ""
        payload = body.encode()
        self.send_response(status)
        self.send_header("Content-Type", mime)
        self.send_header("Content-Length", str(len(payload)))
        for key, value in headers.items():
            self.send_header(key, value)
        self.end_headers()
        if self.command != "HEAD":
            self.wfile.write(payload)


class SmokeTest(unittest.TestCase):
    def test_operator_script_pass_and_fail_without_response_leak(self):
        server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), Fixture)
        thread = threading.Thread(target=server.serve_forever)
        thread.start()
        try:
            command = ["python3", str(Path(__file__).with_name("seo-smoke.py")),
                       "--origin", BASE, "--target", f"http://127.0.0.1:{server.server_port}",
                       "--slug", "public-fixture", "--id", "1", "--hidden-slug", PRIVATE, "--host-check"]
            result = subprocess.run(command, capture_output=True, text=True, check=False)
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertIn("PASS:", result.stdout)
            Fixture.bad_canonical = True
            result = subprocess.run(command, capture_output=True, text=True, check=False)
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("FAIL:", result.stderr)
            self.assertNotIn(PRIVATE, result.stdout + result.stderr)
            self.assertNotIn("incorrect.example", result.stdout + result.stderr)
        finally:
            Fixture.bad_canonical = False
            server.shutdown()
            thread.join()
            server.server_close()


if __name__ == "__main__":
    unittest.main()
