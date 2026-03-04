# build lambda 
build:
	GOOS=linux GOARCH=amd64 go build -o main main.go
# zip lambda
zip:	zip main.zip main
# clean up
clean:	rm -f main main.zip