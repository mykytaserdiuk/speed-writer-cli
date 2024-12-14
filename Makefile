
build-win:
	del "./main.exe"
	go build ./cmd/main.go
	./main.exe

build-lin:
	go build ./cmd/main.go
