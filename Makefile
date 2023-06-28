# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/cli: build the cmd/cli application
.PHONY: build/cli
build/cli:
	cd backend; \
	go mod verify; \
	go build -ldflags='-s' -o=../bin/cli ./cmd/cli