.PHONY: compile-proto
compile-proto:
	@echo "==> Checking buf dependencies..."
ifeq (, $(shell command -v buf 2> /dev/null))
	@echo "==> Setup: Buf not installed, please follow the instructions on https://docs.buf.build/installation"
endif
	@echo "==> compiling proto..."
	@echo "===> generating grpc code..."
	@buf generate
