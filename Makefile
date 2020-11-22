.DEFAULT_GOAL := help
.PHONY: help
help: ## helpを表示
	@echo '  see: https://git.dmm.com/dmm-app/pointclub-api'
	@echo ''
	@grep -E '^[%/0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build
build: ## build
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/butler ./main.go

.PHONY: deploy
deploy: ## deploy
	make build
	scp ./bin/butler hobigon:~/jobs
