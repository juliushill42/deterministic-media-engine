#!/usr/bin/env bash
set -Eeuo pipefail
mkdir -p dist
./dist/dme inspect -in examples/demo.json > dist/report.json
./dist/dme render -in examples/demo.json -out dist/render.wav >/dev/null
[[ -s dist/render.wav ]]
grep -q 'pcm_sha256' dist/report.json
printf 'SMOKE PASS  wav=%s bytes\n' "$(wc -c < dist/render.wav)"
