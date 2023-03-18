SHELL := /bin/bash

build:
	go build -o main main.go

govendor:
	go mod tidy -compat=1.20
	go mod vendor
	git add vendor

cleanup:
	brew uninstall go || echo "failed"
	sudo rm -rf /usr/local/go || echo "failed"
	sudo rm /etc/paths.d/go || echo "failed"
	sudo rm -rf /opt/homebrew/Cellar/go* || echo "failed"

setup:
	HOMEBREW_NO_AUTO_UPDATE=1 brew install go@1.20
