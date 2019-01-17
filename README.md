## build

```shell
docker build -t watchdog-symlinker .
```

## get binary on local

```shell
go get -u github.com/golang/dep/cmd/dep
dep ensure
go build
```