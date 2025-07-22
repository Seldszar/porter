all: binary-windows-amd64 binary-linux-amd64 binary-linux-arm64

binary-windows-amd64:
	CGO_ENABLED=0 GOGC=off GOOS=windows GOARCH=amd64 go build -o "./dist/porter-windows-amd64.exe" .

binary-linux-amd64:
	CGO_ENABLED=0 GOGC=off GOOS=linux GOARCH=amd64 go build -o "./dist/porter-linux-amd64" .

binary-linux-arm64:
	CGO_ENABLED=0 GOGC=off GOOS=linux GOARCH=arm64 go build -o "./dist/porter-linux-arm64" .
