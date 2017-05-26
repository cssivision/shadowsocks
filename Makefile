build:
		docker build -t shadowsocks:personal .
test:
		go test -v -race
cover:
	rm -f *.coverprofile
	go test -coverprofile=shadowsocks.coverprofile
	go tool cover -html=shadowsocks.coverprofile
	rm -f *.coverprofile