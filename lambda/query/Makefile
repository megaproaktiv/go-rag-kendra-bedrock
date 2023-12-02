build-DemoFunction:
	env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap
	cp ./bootstrap $(ARTIFACTS_DIR)/.
