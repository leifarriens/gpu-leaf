## Motivation

When capping fps to a certain amount the gpu uses more power (and produces more heat) than needed. By dynamically reducing the power limit we can reach the same desired fps but with less power usage.

## Utils

List all valid properties to query gpu:

```sh
nvidia-smi --help-query-gpu
```
