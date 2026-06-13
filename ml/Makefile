.PHONY: help data collect clean split train export benchmark infer costs flywheel train-locale export-locale train-all-locales export-all-locales

LOCALES ?= zh-Hans,en-US

help:
	@echo "Targets: data, train, train-all-locales, export-all-locales, benchmark, flywheel"

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

train-locale:
	@test -n "$(LOCALE)" || (echo "Usage: make train-locale LOCALE=zh-Hans" && exit 1)
	python3 train/train_bayes.py --locale $(LOCALE)
	python3 train/train_coreml.py --locale $(LOCALE)
	python3 export/export_bayes_ios.py --locale $(LOCALE)

train-all-locales:
	@for loc in $$(echo "$(LOCALES)" | tr ',' ' '); do \
		$(MAKE) train-locale LOCALE=$$loc; \
	done

export:
	python3 export/export_coreml.py
	python3 export/export_bayes_ios.py

export-locale:
	@test -n "$(LOCALE)" || (echo "Usage: make export-locale LOCALE=zh-Hans" && exit 1)
	python3 export/export_coreml.py --locale $(LOCALE)
	python3 export/export_bayes_ios.py --locale $(LOCALE)

export-all-locales:
	@for loc in $$(echo "$(LOCALES)" | tr ',' ' '); do \
		$(MAKE) export-locale LOCALE=$$loc; \
	done
	CI=true python3 scripts/verify_coreml_export.py

benchmark: train-all-locales export-all-locales
	python3 benchmark/run_benchmark.py
	python3 benchmark/compare_layers.py

infer:
	python3 infer/infer_bayes.py "$(TEXT)"

costs:
	python3 costs/estimate.py

flywheel:
	python3 flywheel/ingest_feedback.py
	python3 flywheel/dedupe_and_merge.py

retrain-all:
	$(MAKE) flywheel
	$(MAKE) train-all-locales
	$(MAKE) export-all-locales
	$(MAKE) benchmark
	MODEL_LOCALES=$(LOCALES) python3 flywheel/publish_to_backend.py
