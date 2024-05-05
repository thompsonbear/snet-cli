BINARY_NAME=snet
build:
	GOARCH=amd64 GOOS=darwin go build -o ./build/macos/snet ./snet
	GOARCH=amd64 GOOS=windows go build -o ./build/windows/snet.exe ./snet
	GOARCH=amd64 GOOS=linux go build -o ./build/linux/snet ./snet

archive: | archive_folder
	tar -czf ./archive/snet-linux-amd64.tar.gz -C ./build/linux .
	tar -czf ./archive/snet-macos-amd64.tar.gz -C ./build/macos .
	tar -czf ./archive/snet-windows-amd64.tar.gz -C ./build/windows .

archive_folder:
	mkdir -p ./archive

clean:
	go clean
	rm -rf ./build
	rm -rf ./archive

