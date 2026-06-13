.PHONY: ios-test gateway backend-dev ml-train ml-benchmark docker-up tier1 tier1-mongodb tier3-mongodb tier4-mtls verify

verify:
	bash scripts/verify.sh

tier1-mongodb:
	bash deploy/tiers/tier1-mongodb.sh

tier3-mongodb:
	bash deploy/tiers/tier3-helm-mongodb.sh

tier4-mtls:
	bash deploy/tiers/tier4-mtls-compose.sh

tier4-helm-mtls:
	bash deploy/tiers/tier4-helm-prod-mtls.sh

ios-test:
	cd apps/ios/Packages/FilterEngine && swift test

gateway:
	cd services/gateway && CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server

backend-dev:
	cd services/rules && go run ./cmd/server &
	cd services/model && PORT=8083 go run ./cmd/server &
	cd services/gateway && CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server

ml-train:
	cd ml && make train

ml-benchmark:
	cd ml && make benchmark

docker-up:
	./deploy/tiers/tier1-compose.sh

tier1:
	./deploy/tiers/tier1-compose.sh
