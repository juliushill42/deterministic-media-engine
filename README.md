# Deterministic Media Engine

A local-first media execution engine where a normalized project graph produces a reproducible render and SHA-256 proof. The shipped native engine renders PCM16 WAV from deterministic tone clips with timeline placement, gain, fades, mixing, clipping, canonical project hashing, HTTP rendering, and an operational timeline UI.

## Run

```bash
make all
./dist/dme serve
# http://127.0.0.1:8789
```

Render without the web UI:

```bash
./dist/dme render -in examples/demo.json -out dist/render.wav
./dist/dme inspect -in examples/demo.json
```

The core has no network dependency and no third-party Go packages. Project validation rejects invalid timing, channels, rates, fades, duplicate tracks, and unsupported sources before execution.

Owner: Julius Cameron Hill / Titan Universal AI LLC  
Watermark: `":"`
