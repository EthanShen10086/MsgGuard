#!/bin/bash
# Schedule retrain — run from cron or K8s CronJob
set -euo pipefail
cd "$(dirname "$0")/.."
make flywheel
make train
make benchmark
python3 flywheel/publish_to_backend.py
