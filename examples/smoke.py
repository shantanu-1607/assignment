#!/usr/bin/env python3
"""Functional local checks of the built reference server; never contacts Cynbit."""
import json
import os
from pathlib import Path
import subprocess
import time
import urllib.error
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
BASE = "http://127.0.0.1:18081"
# Local checks must not send loopback requests through an inherited outbound proxy.
CLIENT = urllib.request.build_opener(urllib.request.ProxyHandler({}))

def call(path, body=None, headers=None, method=None):
    data = json.dumps(body).encode() if body is not None else None
    request = urllib.request.Request(BASE + path, data=data, headers=headers or {}, method=method)
    try:
        response = CLIENT.open(request, timeout=2)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        return response.status, response.headers, response.read()

with subprocess.Popen([str(ROOT / "bin/server")], cwd=ROOT,
                      env={**os.environ, "ADDR": "127.0.0.1:18081", "RATE_PER_MINUTE": "100"},
                      stdout=subprocess.DEVNULL) as process:
    try:
        for attempt in range(50):
            if process.poll() is not None:
                raise RuntimeError("server exited before readiness (port may already be in use)")
            try:
                if call("/readyz")[0] == 200:
                    break
            except (OSError, urllib.error.URLError):
                pass
            time.sleep(0.1)
        else:
            raise RuntimeError("server never became ready")
        status, headers, body = call("/v1/services")
        assert status == 200 and len(json.loads(body)["data"]) == 2
        assert call("/v1/services", headers={"If-None-Match": headers["ETag"]})[0] == 304
        status, _, body = call("/v1/estimates", {"service_id": "api-review", "integrations": 2}, {"Content-Type": "application/json"})
        assert status == 200 and json.loads(body)["data"]["effort_units"] == 7
        status, _, body = call("/v1/estimates", {"service_id": "api-review", "integrations": 21}, {"Content-Type": "application/json"})
        assert status == 422 and json.loads(body)["error"]["code"] == "validation_failed"
        assert call("/missing")[0] == 404
        print("PASS: readiness, catalog, ETag/304, estimate, validation, JSON 404")
    finally:
        process.terminate()
        try:
            code = process.wait(timeout=12)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait()
            raise RuntimeError("graceful shutdown timed out")
        if code != 0:
            raise RuntimeError(f"server shutdown exit status: {code}")
        print("PASS: SIGTERM graceful shutdown, exit status 0")
