#!/usr/bin/env python3
"""Read-only public HTTP acceptance with curl; never print response bodies/secrets."""

import argparse
import html.parser
import json
import re
import subprocess
import tempfile
import urllib.parse
import xml.etree.ElementTree as ET
from pathlib import Path


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


def origin(raw):
    value = urllib.parse.urlsplit(raw)
    require(value.scheme in ("http", "https") and value.hostname and
            not value.username and not value.password and value.path in ("", "/") and
            "?" not in raw and "#" not in raw, "invalid explicit origin")
    return raw.rstrip("/")


class Head(html.parser.HTMLParser):
    def __init__(self):
        super().__init__()
        self.values = {}
        self.ld = []
        self.in_ld = False

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        key = attrs.get("property", attrs.get("name")) if tag == "meta" else None
        value = attrs.get("content")
        if tag == "link":
            key, value = attrs.get("rel"), attrs.get("href")
        if key:
            self.values.setdefault(key, []).append(value)
        if tag == "script" and attrs.get("type") == "application/ld+json":
            self.in_ld = True

    def handle_data(self, data):
        if self.in_ld:
            self.ld.append(json.loads(data))

    def handle_endtag(self, tag):
        if tag == "script":
            self.in_ld = False


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--origin", required=True, help="configured PUBLIC_SITE_URL")
    parser.add_argument("--target", help="optional direct backend origin; does not change expected canonical")
    parser.add_argument("--slug", required=True, help="operator-owned public fixture slug")
    parser.add_argument("--id", required=True, type=int, help="same public article numeric ID")
    parser.add_argument("--hidden-slug", required=True, help="operator-owned unpublished fixture; not printed")
    parser.add_argument("--host-check", action="store_true", help="send hostile Host to direct Go backend (not a virtual-host ingress)")
    args = parser.parse_args()
    canonical = origin(args.origin)
    target = origin(args.target or canonical)
    for slug in (args.slug, args.hidden_slug):
        require(re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", slug) and
                not slug.isdigit() and len(slug) <= 120, "invalid fixture slug")
    require(args.id > 0, "invalid fixture ID")

    with tempfile.TemporaryDirectory(prefix="inkspace-seo-smoke-") as directory:
        root = Path(directory)

        def request(path, method="GET", etag=None):
            headers, body = root / "headers", root / "body"
            # curl may not create/truncate --output for a bodyless 304; never
            # mistake bytes from the preceding GET for this response's body.
            body.write_bytes(b"")
            command = ["curl", "--disable", "--silent", "--show-error", "--max-time", "20",
                       "--proto", "=http,https", "--dump-header", str(headers),
                       "--output", str(body), "--write-out", "%{http_code}"]
            if method == "HEAD":
                command.append("--head")
            if etag:
                command.extend(["--header", "If-None-Match: " + etag])
            if args.host_check:
                command.extend(["--header", "Host: hostile.invalid",
                                "--header", "X-Forwarded-Host: forwarded.invalid"])
            # Never follow redirects, send credentials/cookies or print URLs/body.
            result = subprocess.run(command + [target + path], capture_output=True, check=False)
            require(result.returncode == 0, "curl transport failure (response withheld)")
            lines = headers.read_text().splitlines()
            parsed = {}
            for line in lines:
                if line.startswith("HTTP/"):
                    parsed = {}
                elif ":" in line:
                    key, value = line.split(":", 1)
                    parsed[key.lower()] = value.strip()
            payload = b"" if method == "HEAD" else body.read_bytes()
            require(b"hostile.invalid" not in payload and b"forwarded.invalid" not in payload,
                    "untrusted Host leaked")
            return int(result.stdout), parsed, payload

        article_url = canonical + "/blog/" + args.slug
        status, headers, body = request("/blog/" + args.slug)
        require(status == 200 and headers.get("cache-control") == "no-cache", "slug HTML status/cache failed")
        head = Head()
        head.feed(body.decode("utf-8"))
        require(head.values.get("canonical") == [article_url] and
                head.values.get("og:url") == [article_url] and
                len(head.values.get("description", [])) == 1 and
                len(head.values.get("twitter:card", [])) == 1 and
                len(head.ld) == 1 and head.ld[0].get("@type") == "BlogPosting" and
                head.ld[0].get("url") == article_url, "initial metadata mismatch")
        require(head.values.get("alternate") == [canonical + "/feed.xml", canonical + "/atom.xml"],
                "HTML feed alternate mismatch")
        for method in ("GET", "HEAD"):
            status, headers, body = request("/blog/" + str(args.id), method)
            require(status == 301 and headers.get("location") == article_url, "legacy redirect mismatch")
            status, _, body = request("/blog/" + args.hidden_slug + "?status=0&show_all=true", method)
            require(status == 404 and args.hidden_slug.encode() not in body, "hidden article not uniformly 404")
        status, _, _ = request("/blog/" + args.slug, "HEAD")
        require(status == 200, "article HEAD failed")
        status, _, body = request("/api/settings/public")
        settings = json.loads(body)
        require(status == 200 and settings.get("code") == 0 and
                settings["data"].get("public_site_url") == canonical, "runtime canonical config mismatch")

        documents = [("/feed.xml", "application/rss+xml", "rss"),
                     ("/atom.xml", "application/atom+xml", "{http://www.w3.org/2005/Atom}feed"),
                     ("/sitemap.xml", "application/xml", None),
                     ("/sitemaps/1.xml", "application/xml", "{http://www.sitemaps.org/schemas/sitemap/0.9}urlset"),
                     ("/robots.txt", "text/plain", None)]
        for path, mime, tag in documents:
            status, headers, body = request(path)
            require(status == 200 and headers.get("content-type", "").lower() == mime + "; charset=utf-8" and
                    headers.get("cache-control") == "public, max-age=300" and headers.get("etag"),
                    "discovery status/MIME/cache failed: " + path)
            require(("/blog/" + args.hidden_slug).encode() not in body, "hidden fixture leaked to discovery")
            if path == "/robots.txt":
                require(("Sitemap: " + canonical + "/sitemap.xml").encode() in body, "robots sitemap mismatch")
            else:
                xml = ET.fromstring(body)
                if tag:
                    require(xml.tag == tag, "XML namespace/root mismatch: " + path)
                else:
                    require(xml.tag in ("{http://www.sitemaps.org/schemas/sitemap/0.9}urlset",
                                        "{http://www.sitemaps.org/schemas/sitemap/0.9}sitemapindex"), "invalid sitemap root")
                # URL-bearing discovery fields must share the configured origin.
                for element in xml.iter():
                    local = element.tag.split("}")[-1]
                    if local in ("loc", "guid", "id", "link"):
                        url = element.attrib.get("href", element.text or "").strip()
                        if url:
                            require(url.startswith(canonical + "/"), "noncanonical discovery URL")
            etag = headers["etag"]
            status, head_headers, _ = request(path, "HEAD")
            require(status == 200 and head_headers.get("etag") == etag and
                    head_headers.get("content-length") == headers.get("content-length"), "discovery HEAD mismatch")
            for method in ("GET", "HEAD"):
                status, _, payload = request(path, method, "W/" + etag)
                require(status == 304 and not payload, "conditional request failed")
        print("PASS: public HTML/redirect/404/runtime origin, XML/TXT, shard, HEAD and ETag (no bodies logged)")


if __name__ == "__main__":
    try:
        main()
    except (RuntimeError, ValueError, KeyError, ET.ParseError, OSError):
        # Do not echo parser/transport exceptions containing user-controlled data.
        raise SystemExit("FAIL: SEO smoke acceptance; inspect safely in an isolated environment (response withheld)")
