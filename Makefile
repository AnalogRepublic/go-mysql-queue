test:
	cd msq; go test -v -race

link:
	ln -s $(PWD) ~/go/src
