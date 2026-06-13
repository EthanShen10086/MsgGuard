.PHONY: help data collect clean split train export benchmark infer costs flywheel

help:
	@echo "Targets: data, train, export, benchmark, infer, costs, flywheel"

collect:
	python3 pipeline/collect_github.py || true
	python3 pipeline/collect_uci.py || true

clean:
	python3 pipeline/generate_seed.py
	python3 pipeline/clean.py

split:
	python3 pipeline/split.py
	python3 pipeline/label.py
	python3 pipeline/registry.py

data: collect clean split

train:
	python3 train/train_bayes.py
	python3 train/train_coreml.py

export:
	python3 export/export_coreml.py

benchmark: train
	python3 benchmark/run_benchmark.py
	python3 benchmark/compare_layers.py

infer:
	python3 infer/infer_bayes.py "$(TEXT)"

costs:
	python3 costs/estimate.py

flywheel:
	python3 flywheel/ingest_feedback.py
	python3 flywheel/dedupe_and_merge.py
