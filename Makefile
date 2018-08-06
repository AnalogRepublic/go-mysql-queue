test:
	go test -v ./msq

test-race:
	go test -v -race ./msq

link:
	ln -s $(PWD) ~/go/src
