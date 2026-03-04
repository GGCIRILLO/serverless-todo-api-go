# build lambda binary for AWS
build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda

# zip lambda package
zip: build
	zip function.zip bootstrap

# clean up generated files
clean:
	rm -f bootstrap function.zip
