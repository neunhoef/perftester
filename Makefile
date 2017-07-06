all:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" && mv perftester perftester.linux
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" && mv perftester perftester.mac
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" && mv perftester.exe perftester.windows.exe
	
