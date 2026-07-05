# Twitch Channel Points Miner — Go Edition

A high-performance Go rewrite of [Twitch-Channel-Points-Miner-v2](https://github.com/rdavydov/Twitch-Channel-Points-Miner-v2). Mines channel points, claims bonuses, places predictions, joins raids, claims drops — all with a fraction of the resource usage.

| | Python | Go |
|-|--------|----|
| Memory | ~80–120 MB | ~5–15 MB |
| Docker image | ~800 MB | ~10–15 MB |
| Startup | ~5–10 s | < 100 ms |
| Threads | 60+ | ~10–20 goroutines |

## Wiki pages

| Page | What it covers |
|------|----------------|
| [Getting Started](Getting-Started) | From zero to running in under 5 minutes |
| [Configuration Reference](Configuration-Reference) | Every YAML option with type, default and description |
| [Authentication](Authentication) | All 5 auth methods — when and how to use each |
| [Prediction Strategies](Prediction-Strategies) | All 8 strategies explained with examples and tips |
| [Notifications](Notifications) | Setup guide for all 6 providers + batching deep-dive |
| [Troubleshooting](Troubleshooting) | Common errors, their causes and fixes |
| [Architecture](Architecture) | Package map and data flow for contributors |

## Key resources

- [README](https://github.com/Guliveer/twitch-miner-go#readme) — installation, Docker, Fly.io, service setup
- [Telemetry dashboard](https://github.com/Guliveer/twitch-miner-go-telemetry) — anonymous usage data server (open source)
- [configs/example.yaml.example](https://github.com/Guliveer/twitch-miner-go/blob/main/configs/example.yaml.example) — fully annotated config template
- [CONTRIBUTING.md](https://github.com/Guliveer/twitch-miner-go/blob/main/CONTRIBUTING.md) — commit convention, git hooks, automated versioning
- [Releases](https://github.com/Guliveer/twitch-miner-go/releases) — changelog and downloads
- [Discussions](https://github.com/Guliveer/twitch-miner-go/discussions) — questions and ideas
