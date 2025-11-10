# gpu-leaf

Adaptive NVIDIA GPU power limit controller for steady FPS at lower power and noise.

When you cap FPS, the GPU often uses more power (and generates more heat) than necessary. gpu-leaf samples utilization and reduces the power limit up or down to meet demand, helping you maintain the same FPS while saving power and reducing fan noise.

> Note: This tool uses `nvidia-smi` and requires a power-manageable NVIDIA GPU.

## Install

Build locally:

```sh
./scripts/build.sh
```

The binary will be placed at `bin/gpu-leaf` with an embedded version string if Git metadata is available.

## Usage

Basic example (poll every second, threshold 95% on GPU 0, default power ceiling):

```sh
bin/gpu-leaf --interval 1000 --threshold 95 --gpu 0
```

Allow raising the power limit up to the GPU's max (above the default limit):

```sh
bin/gpu-leaf --overclock
```

Dry-run (log intended changes, make no system modifications):

```sh
bin/gpu-leaf --dry-run
```

### Flags

- `--interval, -l` (ms): polling interval. Default: `1000`.
- `--threshold, -t` (%): utilization threshold to decide whether to increase or decrease the power limit. Default: `95`.
- `--gpu` (index): GPU index to control in multi-GPU systems. Default: `0`.
- `--overclock, --oc` (bool): allow using the maximum power limit instead of the default ceiling. Default: `false`.
- `--dry-run` (bool): log intended changes without applying them. Default: `false`.
- `--log-stdout` (bool): log to stdout. Default: `true`.
- `--log-file` (path): also log to a file. Set to empty to disable file logging. Default: `gpu_leaf.log`.
- `--version` (bool): print version and exit.

### How it works (algorithm)

Every interval, gpu-leaf reads GPU temperature, power draw, utilization, and the current power limit via `nvidia-smi`. If utilization is below the threshold, it reduces the power limit proportionally; if above, it increases it, bounded between the minimum power limit and either the default or maximum power limit (when `--overclock` is enabled). Changes are applied in integer watts as required by `nvidia-smi`.

### Safety and permissions

- If you see permission errors when setting the power limit, try running with elevated privileges or adjust system policies to allow `nvidia-smi -pl`.
- Use `--dry-run` first to validate observed behavior before applying.
- `--overclock` allows the power limit to increase above the default limit; use this only if you understand the thermals of your system.

## Utilities

List all valid GPU query properties:

```sh
nvidia-smi --help-query-gpu
```

## Development

- Go version: declared in `go.mod` (currently `1.19`)
- Build script: `scripts/build.sh` (embeds version via ldflags)
- Main entrypoint: `cmd/gpuleaf/main.go`
- Core logic: `internal/gpu`
- Utilities: `internal/utils`
