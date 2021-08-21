# Project Setup Steps

* [install protoc](https://grpc.io/docs/protoc-installation/)

* [install the Go protocol buffers plugin](https://developers.google.com/protocol-buffers/docs/gotutorial)
    * `go install google.golang.org/protobuf/cmd/protoc-gen-go`

 * Set up module and run the program(s)   
    * `export PATH=~/go/bin:$PATH`
    * `protoc -I=./ --go_out=./ ./addressbook.proto`
    * `go mod init protobufdemo`
    * `go tidy`
    * `go run add_person.go addressbook.data`
    * `go run list_person.go addressbook.data`