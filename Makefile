test:
	cd msq; go test -v

test-race:
	cd msq; go test -v -race

link:
	ln -s $(PWD) ~/go/src
