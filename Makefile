.PHONY: build
build:
	tinygo build -o ./bin/main.go.wasm -scheduler=none -target=wasi main.go