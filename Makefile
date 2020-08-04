# Makefile to help with builds
#

VERSION      := 0.0.1
GIT_COMMIT   := $(shell git rev-parse HEAD)
GO_VERSION   := $(shell go version | sed 's/go version //')
BIN_NAME     := goldap

# Use the above vars as build flags
LDFLAGS=-ldflags '\
	-X "main.Version=${VERSION}" \
	-X "main.GitCommit=${GIT_COMMIT}" \
	-X "main.GoVersion=${GO_VERSION}" \
	'

# Colors
CCEND=\033[0m
RED=\033[0;31m
GREEN=\033[0;32m

help:
	@echo
	@echo "--------------------------------------------------------------"
	@echo
	@echo "Helper commands for building and deploying the project"
	@echo
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "--------------------------------------------------------------"
	@echo

clean: ## Clean local binaries and cache
	@# echo "Cleaning local binaries and caches"
	@go clean
	@rm -rf dist

build: ## Build binary for the current platform
	@# echo "Building binary for your current environment"
	@go build ${LDFLAGS}

debug: ## Build binary and exec with delve debugger
	go build ${LDFLAGS} -gcflags="-N -l"
	dlv exec goldap --build-flags="-gcflags='-N -l'" -- search

install: ## Install binary in $GOPATH
	@echo "Installing binary for your current environment"
	go install ${LDFLAGS}

build_all: ## Build multi-platform binaries
	@echo "Building multi-platform binaries using gox"
	gox ${LDFLAGS} -output "dist/${BIN_NAME}_{{.OS}}_{{.Arch}}"

debug_in_tmux_pane: ## Launch the delve debugger in the bottom-right tmux pane
	@echo "Launching delve debugger in bottom-right tmux pane"
	@if pgrep dlv >/dev/null 2>&1; then killall dlv; fi
	@tmux send-keys -t bottom-right 'echo "Magic launching Delve debugger"' Enter
	tmux send-keys -t bottom-right 'dlv exec ${BIN_NAME} -- ${ARGS}' ENTER

exec_in_tmux_pane:
	@echo "Executing in bottom-right tmux pane"
	@if pgrep dlv >/dev/null 2>&1; then killall dlv; fi
	@tmux send-keys -t bottom-right './${BIN_NAME} ${ARGS}' ENTER

execloop:
	@echo "Starting file watcher"
	@fswatch --exclude='.*\.git'\
   	--exclude='.*\.yaml'\
   	--exclude='.*\.json'\
	--exclude='.*\.swp'\
	--exclude='*.md'\
	--exclude='.*debug.*?'\
	--exclude='.*4913'\
	--exclude='fixtures'\
	--exclude='LICENSE'\
	--exclude='Makefile'\
	--exclude='${BIN_NAME}'\
	--recursive . |\
	xargs -n1 -I{} sh -c 'echo "Change detected: {}"; make clean build exec_in_tmux_pane ARGS="${ARGS}";'

debugloop:
	@echo "Starting file watcher"
	fswatch --exclude='.*' --include='.*\.go$$' --recursive . | xargs -n1 -I{} sh -c 'echo; echo "Change detected: {}"; make clean build debug_in_tmux_pane ARGS="${ARGS}";'
