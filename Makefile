DEFAULT_TARGET: run

.PHONY: run-linux
run-linux:
	@PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig go run main.go

.PHONY: run
run:
	@go run main.go
