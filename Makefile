.PHONY: ;
.SILENT: ;               # no need for @
.ONESHELL: ;             # recipes execute in same shell
.NOTPARALLEL: ;          # wait for target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

EXEC=minio-uploader
OUTPUT=output

build: clean build-win build-mac

build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "${OUTPUT}/${EXEC}.exe" .

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o "${OUTPUT}/${EXEC}" .

clean:
	rm -rf ${OUTPUT}
