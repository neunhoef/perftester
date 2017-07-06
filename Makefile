all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" && mv perftester perftester.linux
	GOOS=darwin GOWARCH=amd64 go build -ldflags "-s -w" && mv perftester perftester.mac
	
