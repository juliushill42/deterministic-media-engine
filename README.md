# Deterministic Media Engine

A local-first Go media execution engine where a normalized project graph produces a reproducible render plus a SHA-256 proof.

The shipped native engine renders PCM16 WAV from deterministic tone clips.

## Implemented

- Timeline placement.
- Gain.
- Fades.
- Track mixing.
- Clipping.
- Canonical project hashing.
- HTTP rendering.
- Operational timeline UI.
- Input validation before execution.
- No third-party Go packages in the core build path.
- No network dependency for the core renderer.

Validation rejects invalid timing, channels, sample rates, fades, duplicate tracks, and unsupported sources.

## Build and verify

```bash
make all
```

`make all` runs:

```text
go test ./...
go vet ./...
CGO_ENABLED=0 go build ...
./scripts/smoke.sh
```

The resulting binary is:

```text
dist/dme
```

## Run the local service

```bash
./dist/dme serve
```

Default UI/service:

```text
http://127.0.0.1:8789
```

## Render without the UI

```bash
./dist/dme render \
  -in examples/demo.json \
  -out dist/render.wav
```

Inspect the normalized project:

```bash
./dist/dme inspect -in examples/demo.json
```

## One-shot bootstrap

```bash
./bootstrap.sh
```

The bootstrap runs tests, vet, a stripped static-style Go build with `CGO_ENABLED=0`, and the smoke test.

## Repository layout

```text
cmd/dme/       CLI and service entrypoint
internal/      media engine implementation
contracts/     project contracts
examples/      deterministic example projects
scripts/       smoke verification
dist/          binary and rendered output
bootstrap.sh   one-shot build + verify
```

## Current boundary

The native renderer in this repository is centered on deterministic PCM16 WAV/tone-clip execution. The README does not claim codecs or media operations that are not present in this build.

## Ownership

Owner: Julius Cameron Hill / Titan Universal AI LLC  
Watermark: `":"`
