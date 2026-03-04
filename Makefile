# build lambda 
build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda
# zip lambda
zip:	zip function.zip bootstrap
# clean up
clean:	rm -f main main.zip