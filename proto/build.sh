PATH="$PATH:${GOPATH}/bin:${HOME}/go/bin" protoc --go_out=../failureDetection/ ./*.proto