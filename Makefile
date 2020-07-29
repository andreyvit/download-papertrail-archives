
release-mac-intel:
	mkdir -p _dist
	GOOS=darwin GOARCH=amd64 go build -o _dist/download-papertrail-archives-mac-intel
# release-mac-arm64:
# 	mkdir -p _dist
# 	GOOS=darwin GOARCH=arm64 go build -o _dist/download-papertrail-archives-mac-arm64

release-windows-intel:
	mkdir -p _dist
	GOOS=windows GOARCH=amd64 go build -o _dist/download-papertrail-archives-windows-intel
release-windows-arm:
	mkdir -p _dist
	GOOS=windows GOARCH=arm go build -o _dist/download-papertrail-archives-windows-arm

release-linux-intel:
	mkdir -p _dist
	GOOS=linux GOARCH=amd64 go build -o _dist/download-papertrail-archives-linux-intel
release-linux-arm64:
	mkdir -p _dist
	GOOS=linux GOARCH=arm64 go build -o _dist/download-papertrail-archives-linux-arm64

release: release-mac-intel release-windows-intel release-windows-arm release-linux-intel release-linux-arm64
