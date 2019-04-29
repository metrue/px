clean:
	rm -rf px
mod:
	go mod download
build: mod
	go build  -o px px.go
	go build  -o px-daemon ./daemon/daemon.go
