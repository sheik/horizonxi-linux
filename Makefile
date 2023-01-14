
build:
	CGO_ENABLED=0 go build ./cmd/horizonxi-installer
	upx horizonxi-installer
