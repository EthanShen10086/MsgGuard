#!/bin/bash
# Schedule per-locale retrain — run from cron or K8s CronJob
set -euo pipefail
cd "$(dirname "$0")/.."
export MODEL_LOCALES="${MODEL_LOCALES:-zh-Hans,en-US}"
make retrain-all
